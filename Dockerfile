FROM golang:1.19.7-alpine

COPY . /app

WORKDIR /app

RUN go mod download

RUN mkdir -p ./bin
RUN go build -o /bin/tgavatar /app/cmd/

EXPOSE 8081

CMD [ "/bin/tgavatar" ]
