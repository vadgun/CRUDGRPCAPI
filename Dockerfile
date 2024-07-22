FROM golang:1.23rc2-alpine3.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080
EXPOSE 50051

RUN go build -o myapp .

CMD ["./myapp"]
