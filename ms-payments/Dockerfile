FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN apk update && \
    apk add curl
    
RUN go build -o ./main
EXPOSE 8080

CMD [ "./main" ]