# goバージョン
FROM golang:1.19.3-alpine

# アップデートとgitのインストール
RUN apk update && apk add git
# go-websocketディレクトリの作成
RUN mkdir ./go
# ワーキングディレクトリの設定
WORKDIR /go/src/github.com/go-websocket
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . ./go/src/github.com/go-websocket
# パッケージのインポート

RUN go mod init example.com/go-websocket

RUN go get golang.org/x/tools/cmd/goimports
# 追加です
RUN go get github.com/gorilla/websocket