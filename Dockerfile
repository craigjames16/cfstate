FROM golang:latest

WORKDIR /go/src/github.com/craigjames16/cfstate
COPY . .

RUN go get ./...
RUN go install .