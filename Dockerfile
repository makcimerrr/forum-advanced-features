FROM golang:latest
LABEL desc="This is my website"
WORKDIR /app
ADD . /app
EXPOSE 8080
RUN go build -o main cmd/main.go
CMD [ "./main" ]