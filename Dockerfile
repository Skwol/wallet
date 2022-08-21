ARG GOBIN=/app

FROM golang:1.18-alpine as builder
ARG GOBIN
COPY . /go/src/github.com/skwol/wallet
WORKDIR /go/src/github.com/skwol/wallet
RUN apk add --update make gcc
RUN GOBIN=$GOBIN
RUN make install-tools
RUN make build

# Dev image
FROM builder AS dev
EXPOSE 8080
ENTRYPOINT make watch

FROM alpine:3.9
ARG GOBIN
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder ${GOBIN}/wallet wallet
COPY --from=builder ${GOBIN}/walletctl walletctl
EXPOSE 8080 8080
CMD ["./wallet"]