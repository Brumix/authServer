FROM golang:alpine

LABEL maintainer="Bruno Pereira" \
      version="1.0"
#GRPC
ENV GRPCPORT=9000

#Logs
ENV LOGFILE=""
ENV LOGLEVEL="error"

#DataBase
ENV HOST=database
ENV DBUSER=postgres
ENV PASSWORD=pass
ENV DBNAME=hypermediastockapi
ENV DBPORT=5432

#Auth
ENV JWTKEY="2dafede0ed481426c936de1e82173c9e9b58fc22c95bf0cb97862be2bc79daf92092edd0adc2c6b99c4cdf824a8865154f7a892de565b1efb14607d6dd59e836"


WORKDIR /go/src

COPY . .

RUN go mod tidy

EXPOSE $GRPCPORT

RUN go build -o app .


ENTRYPOINT ["./app"]
