# goバージョン
FROM golang:1.19.3-alpine

# アップデートとgitのインストール
RUN apk update && apk add git && apk add gcc && apk add musl-dev
# go-websocketディレクトリの作成
RUN mkdir ./go
# ワーキングディレクトリの設定
WORKDIR /go/src/github.com/go-websocket
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . ./go/src/github.com/go-websocket
