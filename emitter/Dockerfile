FROM golang:1.15

RUN apt-get update || exit 0
RUN apt-get upgrade -y

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main ./server

EXPOSE 8080

CMD ./main