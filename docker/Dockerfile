FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mini-blog ./cmd/mini-blog

CMD ["./mini-blog"]
