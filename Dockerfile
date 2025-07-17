# STEP 1: ビルド用のコンテナ
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先に依存関係をダウンロードしてキャッシュする
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# STEP 2: 実行用の軽量コンテナ
FROM alpine:latest

WORKDIR /app

# ビルド用コンテナから実行ファイルのみをコピー
COPY --from=builder /app/main .

# コンテナ起動時に実行するコマンド
CMD ["./main"]