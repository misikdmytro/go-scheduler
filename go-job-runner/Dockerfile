FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/app ./cmd/app/main.go

WORKDIR /out
CMD ["./app"]