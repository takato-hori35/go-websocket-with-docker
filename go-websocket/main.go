package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

// メッセージ用構造体
type Message struct {
	Msg  string `json:"msg"`
	Name string `json:"name"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

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
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("Cookie: ", err)
	}

	msg.Msg = r.FormValue("msg")
	msg.Name = cookie.Value
	broadcast <- msg
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
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func setName(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("set-name/index.html")
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, nil); err != nil {
		panic(err.Error())
	}
}

func sendCookies(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name: "username",
		Value: r.FormValue("nametxt"),
		Path: "/",
	}

	// cookieに送信
	http.SetCookie(w, cookie)

	// これは結局データをpostしてるだけ
	t := template.Must(template.ParseFiles("static/index.html"))
	if err := t.Execute(w, cookie); err != nil {
		panic(err.Error())
	}
}

func main() {
	portNumber := "9000"
	http.HandleFunc("/ws", websocketConnectHandler)
	http.HandleFunc("/msg", messageHandler)
	http.HandleFunc("/", setName)
	http.HandleFunc("/chat-page", sendCookies) // cookieをsetするapiを用意
	log.Println("Server listening on port ", portNumber)
	go websocketMessages()
	log.Println("success")
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}
