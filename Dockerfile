FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

# -------- Runtime Stage --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .


EXPOSE 8080

CMD ["./server"]
