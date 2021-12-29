FROM golang:1.17 AS builder
ENV CGO_ENABLED=0
COPY go.mod go.sum /go/src/github.com/skwol/wallet/
WORKDIR /go/src/github.com/skwol/wallet
RUN go mod download
COPY . /go/src/github.com/skwol/wallet
RUN GOOS=linux go build -a -installsuffix cgo -o build/app github.com/skwol/wallet/cmd/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/skwol/wallet/build/app ./
COPY --from=builder /go/src/github.com/skwol/wallet/templates ./templates/
EXPOSE 8080 8080
CMD ["./app"]