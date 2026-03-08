FROM golang:latest

RUN apt-get update && apt-get install -y \
    libopencv-dev \
    pkg-config \
    cmake \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY . .
RUN go build -o server .

WORKDIR /app
RUN cp /build/server /app/server && rm -rf /build

EXPOSE 80
EXPOSE 8080

ENTRYPOINT ["./server"]