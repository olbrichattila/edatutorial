# Builder,
# Note this file does not work directly here, as it requires the shared components, it is used by build.sh
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .

WORKDIR /app/build/code
WORKDIR /app/build/code
RUN go build -o app .

# Built
FROM debian:bookworm-slim
COPY --from=builder /app/build/code/app /app/app
CMD ["/app/app"]