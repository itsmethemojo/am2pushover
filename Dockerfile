FROM golang:alpine as builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /root
RUN mkdir -p /root
COPY go.mod go.sum *.go /root/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -installsuffix=cgo am2pushover.go

FROM alpine
RUN apk update && apk add git && apk add ca-certificates

COPY --from=builder /root/am2pushover /root/
EXPOSE 5001
ENTRYPOINT ["/root/am2pushover"]

