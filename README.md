# 簡易チャットシステム
goでwebsocketn追加た学ぶ

## 起動方法
``` console
# コンテナの起動
docker-compose up

# コンテナに入ってbashシェルを起動してコマンドの入力を受け付ける状態にする
docker exec -it go-websocket /bin/sh

# 起動
go run main.go

# アクセス
http://localhost:9000
```