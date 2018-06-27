// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ttd

import "time"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan SendMsg
	multicast chan SendMsg
	// Register requests from the clients.
	register chan *ClientRegister

	// Unregister requests from clients.
	unregister chan *Client
}
type ClientRegister struct {
	client    *Client
	retClient chan *Client
}
// send message struct
type SendMsg struct {
	serverName  string
	sender    	int64
	message    []byte
	recverID  	[]int64
	recvername string
	sendername string
	addTime time.Time
}
func NewHub() *Hub {
	hub := &Hub{
		multicast:  make(chan SendMsg),
		broadcast:  make(chan SendMsg),
		register:   make(chan *ClientRegister),
		unregister: make(chan *Client),
		clients:    make(map[string]map[int64]*Client),
	}
	hub.clients["ttd"] = map[int64]*Client{}
	return hub
}
func (h *Hub) Run() {
	for {
		select {
		case cRegister := <-h.register:
			//h.clients[client] = true
			// 根据传入的client对指定服务器的userid进行socket注册
			if _, ok := h.clients[cRegister.client.serverName]; !ok {
				cRegister.retClient <- cRegister.client
				close(cRegister.client.send)
				continue
			}
			//util.LogInfo().Println(cRegister.client.userID)
			// 判断用户是否已经存在
			if client, ok := h.clients[cRegister.client.serverName][cRegister.client.userID]; ok {
				//重复登录,返回旧的client
				client.repeatLogin = true
				cRegister.retClient <- client
				//用新的客户端覆盖旧的客户端
				h.clients[cRegister.client.serverName][cRegister.client.userID] = cRegister.client
				continue
			}
			h.clients[cRegister.client.serverName][cRegister.client.userID] = cRegister.client
			cRegister.retClient <- cRegister.client

		case client := <-h.unregister:
			//if _, ok := h.clients[client]; ok {
			//	delete(h.clients, client)
			//	close(client.send)
			//}
			//if client.repeatLogin {
			//	close(client.send)
			//	continue
			//}

			// 收到失败的socket就进行注销
			if _, ok := h.clients[client.serverName][client.userID]; ok {
				close(client.send)
				delete(h.clients[client.serverName], client.userID)
			}
		case sendMsg := <-h.multicast:
			// 对指定的用户群发送消息
			for _, userID := range sendMsg.recverID {
				client, ok := h.clients[sendMsg.serverName][userID]
				if !ok {
					continue
				}
				select {
				case client.send <- sendMsg.message:
				default:
					close(client.send)
					delete(h.clients[client.serverName], client.userID)
				}
			}
		case sendMsg := <-h.broadcast:
			//for client := range h.clients {
			//	select {
			//	case client.send <- message.:
			//	default:
			//		close(client.send)
			//		delete(h.clients, client)
			//	}
			//}
			// 对所有的在线用户发送消息
			for userID := range h.clients[sendMsg.serverName] {

				client := h.clients[sendMsg.serverName][userID]
				select {
				case client.send <- sendMsg.message:
				default:
					close(client.send)
					delete(h.clients[client.serverName], client.userID)
				}
			}
		}
	}
}