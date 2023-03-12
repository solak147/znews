FROM golang:latest

RUN mkdir -p /usr/local/go/src/znews
WORKDIR /usr/local/go/src/znews
ADD . /usr/local/go/src/znews

RUN go mod download
RUN go build ./main.go

EXPOSE 8080
CMD ["./main"]