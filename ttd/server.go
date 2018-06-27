package ttd

import ("net/http"
	"yin/xxd/util"
)

const (
	webSocket = "/ws"
)
func InitWs() {
	hub := NewHub()
	go hub.Run()

	// 初始化路由
	http.HandleFunc(webSocket, func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	addr := util.Config.Ip + ":" + util.Config.ChatPort
	util.LogInfo().Println("websocket start,listen addr:", addr, webSocket)

	// 创建服务器
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		util.LogError().Println("websocket server listen err:", err)
		util.Exit("websocket server listen err")
	}
}

