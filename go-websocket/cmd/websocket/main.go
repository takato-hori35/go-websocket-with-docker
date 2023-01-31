package main

import (
	"log"
	"net/http" // httpの標準パッケージ

	"github.com/gorilla/websocket" // go getで追加した websocketの通信プロトコル
)

// メッセージ用構造体
type Message struct {
	Msg  string `json:"msg"` // json形式でやりとりしますよー
	Name string `json:"name"`
}

var clients = make(map[*websocket.Conn]bool) // 接続されるクライアント
var broadcast = make(chan Message)           // メッセージ用ブロードキャストチャネル

// gorilla/websocketパッケージを使ってアップグレード設定（今回はバッファサイズのみ）
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HTTP通信をアップグレード
func websocketConnectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	clients[conn] = true // 接続したクライアントを保存
}

// チャットの発言用のハンドラ
func messageHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	msg.Msg = r.FormValue("msg")
	msg.Name = r.FormValue("name")
	broadcast <- msg // メッセージを受けたらメッセージ用ブロードキャストチャネルに送信
}

// ブロードキャストチャンネルに通知が来たら、全クライアントに向けてメッセージを送信
func websocketMessages() {
	for {
		// チャネルからメッセージを受け取る
		msg := <-broadcast
		// 現在接続しているクライアントすべてにメッセージを送信する
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println(err) //クライアントの接続が切れるとエラー
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	portNumber := "9000"
	http.Handle("/", http.FileServer(http.Dir("static"))) // "/": staticフォルダにある静的ファイルを公開
	http.HandleFunc("/ws", websocketConnectHandler)       // "/ws": Websocket通信接続用
	http.HandleFunc("/msg", messageHandler)               // "/msg": メッセージ送信用API
	log.Println("Server listening on port ", portNumber)
	go websocketMessages()
	http.ListenAndServe(":"+portNumber, nil) // 9000でWebサーバーを開始
}
