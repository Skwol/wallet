FROM golang:1.17.5-alpine3.15 AS builder
COPY . /go/src/github.com/skwol/wallet
WORKDIR /go/src/github.com/skwol/wallet
RUN go mod verify && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/app github.com/skwol/wallet/cmd/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/skwol/wallet/build/app ./
EXPOSE 8080 8080
CMD ["./app"]