# STEP 1: ビルド用のコンテナ
FROM golang:1.24-alpine AS builder

# SSL証明書をインストール (MySQL接続などで必要)
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# 先に依存関係をダウンロードしてキャッシュする
COPY go.mod go.sum ./
RUN go mod download


# コンテナ起動時に実行するコマンド
CMD ["go", "run", "main.go"]