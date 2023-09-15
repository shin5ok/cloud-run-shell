FROM golang:1.21 AS builder

ENV ROOT=/app
WORKDIR $ROOT
COPY . .
WORKDIR $ROOT
RUN GGO_ENABLED=0 GOOS=linux go build -o ./main .

FROM gcr.io/distroless/base-debian11 AS runner
WORKDIR /
COPY --from=builder /app/main ./main
# USER nobody
CMD ["./main"]
