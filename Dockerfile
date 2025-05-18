FROM golang:1.24.3-alpine3.21

WORKDIR /app

COPY ./cmd/agent ./
COPY go.mod ./

RUN go mod download

COPY  . .

RUN go build -o main .

CMD ["./main", "-pk=secretkey"]