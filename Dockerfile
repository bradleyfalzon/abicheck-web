FROM golang:1.6.3-onbuild
RUN go get github.com/bradleyfalzon/abicheck/cmd/abicheck
EXPOSE 80
WORKDIR /go/src/app
CMD app
