FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Tests
# RUN go tests ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o url-shortner ./cmd/app

FROM alpine:latest AS final

WORKDIR /root/

COPY --from=builder /app/url-shortner .

EXPOSE 8080

CMD [ "./url-shortner" ]