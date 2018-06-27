// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttd

import (
	"net/http"
	"time"
	"yin/xxd/util"
	"github.com/gorilla/websocket"
	"yin/xxd/api"
	"yin/xxd/models"
	"yin/xxd/snowflake"
	"yin/xxd/orm"
	"strconv"
	"github.com/json-iterator/go"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	repeatLogin bool
	// The websocket connection.
	conn *websocket.Conn
	serverName  string          // User server
	// Buffered channel of outbound messages.
	send chan []byte
	userID      int64
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				util.LogError().Printf("error: %v", err)
			}
			util.LogError().Printf("error: %v", err)
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//返回user id 、登录响应的数据、ok
		if dataProcessing(message, c) != nil {

			util.LogInfo().Println("client exit ip:", c.conn.RemoteAddr())
			break
		}
	}
}
func switchMethod(message []byte, parseData api.ParseData, client *Client) error {

	castStr:=parseData.Module() + "." + parseData.Method()
	switch castStr {
	case "chat.login":
		if err := chatLogin(parseData, client); err != nil {
			return err
		}

		break

	case "chat.logout":
		client.conn.Close()
		/*
		   if err := chatLogout(parseData.UserID(), client); err != nil {
			   return err
		   }
		*/
		break

	default:
		err := transitData(message, parseData, client)
		if err != nil {
			util.LogError().Println(err)
		}
		break
	}

	return nil
}
func Online(userid int64,client *Client)(bool)  {
	// 判断用户是否已经存在
	if _, ok := client.hub.clients[client.serverName][userid]; ok {
		return true
	}
	return false
}
//交换数据
func transitData(message []byte,ParseData api.ParseData, client *Client) error {
	userID:=ParseData.UserID()
	if client.userID != userID {
		return util.Errorf("%s", "user id err")
	}
	parseData, err :=api.ApiParse(message,[]byte("token"))
	if err != nil {
		util.LogError().Println("api parse error:", err)
	}
	//添加消息记录
	model:=models.IMRecord{}
	model.Id=snowflake.GetSnowflakeId()
	model.Sendername=api.GetUser(userID).NickName
	model.Sender=strconv.FormatInt(userID,10)
	model.Recver=strconv.FormatInt(ParseData.SendUsers()[0],10)
	model.Recvername=ParseData.Recvername()
	model.Time=time.Now().Format("2006-01-02 15:04:05")
	model.Content=parseData.Content()
	model.Avatar="http://localhost:11443/avatar?id="+strconv.FormatInt(userID,10)
	id,err:=strconv.ParseInt(model.Recver, 10, 64)
	if Online(id,client){
	model.Offline=1//在线
	X2cSend(client.serverName, parseData.SendUsers(), message, client)
	}else {
		model.Offline=0//离线
	}
	orm.AddRecord(&model)
	//x2cMessage, sendUsers, err := api.TransitData(message, client.serverName)
	//if err != nil {
	//    // 与然之服务器交互失败后，生成error并返回到客户端
	//    errMsg, retErr := api.RetErrorMsg("0", "time out")
	//    if retErr != nil {
	//        return retErr
	//    }
	//
	//    client.send <- errMsg
	//    return err
	//}
	return nil
}
// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//会话退出
func chatLogout(userID int64, client *Client) error {
	if client.userID != userID {
		return util.Errorf("%s", "user id error.")
	}
	if client.repeatLogin {
		return nil
	}
	x2cMessage, sendUsers, err := api.ChatLogout(client.serverName, client.userID)
	if err != nil {
		return err
	}
	//util.DelUid(client.serverName,util.Int642String(client.userID))
	return X2cSend(client.serverName, sendUsers, x2cMessage, client)
}
// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.LogError().Println("serve ws upgrader error:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), repeatLogin: false}
	//client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	client.readPump()
}
func X2cSend(serverName string, sendUsers []int64, message []byte, client *Client) error {
	if len(sendUsers) == 0 {
		//send all
		client.hub.broadcast <- SendMsg{serverName: serverName, message: message}
		return nil
	}

	//send users
	client.hub.multicast <- SendMsg{serverName: serverName,recverID: sendUsers, message: message}
	return nil
}
//解析数据.
func dataProcessing(message []byte, client *Client) error {
	parseData, err := api.ApiParse(message, util.Token)
	if err != nil {
		util.LogError().Println("recve client message error")
		return err
	}

	//if util.IsTest && parseData.Test() {
	//    return testSwitchMethod(message, parseData, client)
	//}

	return switchMethod(message, parseData, client)
}
//用户登录
func chatLogin(parseData api.ParseData, client *Client) error {
	user:= api.ChatLogin(parseData)
	//if userID == -1 {
	//    util.LogError().Println("chat login error")
	//    return util.Errorf("%s", "chat login error")
	//}
	//
	//if !ok {
	//    // 登录失败返回错误信息
	//    client.send <- loginData
	//    return util.Errorf("%s", "chat login error")
	//}
	//// 成功后返回login数据给客户端
	//client.send <- loginData
	client.userID = parseData.UserID()
	client.serverName = parseData.ServerName()
	if client.serverName == "" {
		client.serverName = util.Config.DefaultServer
	}

	// 生成并存储文件会员
	//userFileSessionID , err := api.UserFileSessionID(client.serverName, client.userID)
	//if err != nil {
	//    util.LogError().Println("chat user create file session error:", err)
	//    //返回给客户端登录失败的错误信息
	//    return err
	//}
	//// 成功后返回userFileSessionID数据给客户端
	//client.send <- userFileSessionID

	// 获取所有用户列表
	//usergl, err := api.UserGetlist(client.serverName, client.userID)
	//if err != nil {
	//    util.LogError().Println("chat user get user list error:", err)
	//    //返回给客户端登录失败的错误信息
	//    return err
	//}
	//// 成功后返回usergl数据给客户端
	//client.send <- usergl

	// 获取当前登录用户所有会话数据,组合好的数据放入send发送队列
	//getlist, err := api.Getlist(client.serverName, client.userID)
	//if err != nil {
	//    util.LogError().Println("chat get list error:", err)
	//    // 返回给客户端登录失败的错误信息
	//    return err
	//}
	//// 成功后返回gl数据给客户端
	//client.send <- getlist

	//获取离线消息发送给客户端
	offlineMessages,isMsg:= api.GetofflineMessages(client.userID)
	if isMsg{
		client.send <- offlineMessages
	}
	//for _, v := range offlineMessages {
	//	client.send <- v
	//}
	// 推送当前登录用户信息给其他在线用户
	// 因为是broadcast类型，所以不需要初始化userID
	//client.hub.broadcast <- SendMsg{serverName: client.serverName, message: loginData}
	if user.Service==""{
		user.Service="0"
	}
	msg:=models.MessageResult{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if user!=nil{
		msg.Module="chat"
		msg.Method="kickoff"
		msg.Service=user.Service
		msg.Message="登陆成功"
		jsonMsg,err:=json.Marshal(&msg)
		util.Check(err)
		client.send <-jsonMsg
	}else
	{
		msg.Module="chat"
		msg.Method="off"
		msg.Service=user.Service
		msg.Message="登陆失败"
		jsonMsg,err:=json.Marshal(&msg)
		util.Check(err)
		client.send <-jsonMsg
	}
	cRegister := &ClientRegister{client: client, retClient: make(chan *Client)}
	defer close(cRegister.retClient)

	// 以上成功后把socket加入到管理
	client.hub.register <- cRegister
	if retClient := <-cRegister.retClient; retClient.repeatLogin {
		//客户端收到信息后需要关闭socket连接，否则连接不会断开
		retClient.send <- api.RepeatLogin()

		//是重复登录，不需要再发送给其他用户上线信息
		return nil
	}

	return nil
}