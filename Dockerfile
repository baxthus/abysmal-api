FROM golang:alpine as builder

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

# this is not recommended for CI/CD, a better way is to use the install script
RUN go install github.com/go-task/task/v3/cmd/task@latest

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN task build

FROM alpine:latest

COPY --from=builder /build/server /server
COPY --from=builder /build/public /public

CMD ["/server"]