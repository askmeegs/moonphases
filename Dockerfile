
FROM alpine
RUN apk add --update ca-certificates && \
      rm -rf /var/cache/apk/* /tmp/*

FROM golang:1.11-alpine3.7 as build
COPY . $GOPATH/src/github.com/m-okeefe/moonphases
WORKDIR $GOPATH/src/github.com/m-okeefe/moonphases
RUN go build -o ./bin/moonphases ./
ENTRYPOINT ["./bin/moonphases"]
EXPOSE 8001

