FROM golang:1.21 AS builder

ENV ROOT=/app
WORKDIR $ROOT
COPY . .
WORKDIR $ROOT
RUN GGO_ENABLED=0 GOOS=linux go build -o ./main .

FROM ubuntu AS runner
WORKDIR /
COPY --from=builder /app/main ./main
RUN apt-get update && \
    apt-get -y install procps iproute2 net-tools
# USER nobody
CMD ["./main"]
