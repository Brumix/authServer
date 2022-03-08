FROM golang:alpine

LABEL maintainer="Bruno Pereira" \
      version="1.0"

WORKDIR /go/src

COPY . .

RUN go mod tidy

EXPOSE $GRPCPORT

RUN go build -o app .


ENTRYPOINT ["./app"]
