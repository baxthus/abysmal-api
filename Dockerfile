FROM golang:1.20-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN task build

COPY --from=builder /build/server /server

CMD ["/server"]