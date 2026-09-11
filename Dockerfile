FROM golang:1.26
WORKDIR /app

COPY go.mod go.sum /app
RUN go mod download 

COPY . .

RUN go build -o myapp .
CMD ["./myapp"]