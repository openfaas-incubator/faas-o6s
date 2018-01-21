FROM golang:1.9

RUN mkdir -p /go/src/github.com/openfaas-incubator/faas-o6s/

WORKDIR /go/src/github.com/openfaas-incubator/faas-o6s

COPY . .

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faas-o6s .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=0 /go/src/github.com/openfaas-incubator/faas-o6s/faas-o6s .

ENTRYPOINT ["./faas-o6s"]
CMD ["-logtostderr"]
