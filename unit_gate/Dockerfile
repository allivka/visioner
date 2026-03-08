FROM golang:latest AS builder

WORKDIR /build
COPY . .

RUN go build -o server .

FROM debian:bookworm-slim AS runner

WORKDIR /app
COPY --from=builder /build/server /app/server

EXPOSE 80
EXPOSE 8080

ENTRYPOINT [ "./server"]
CMD [""]