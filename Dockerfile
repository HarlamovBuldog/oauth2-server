FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN GOOS=linux CGO_ENABLED=0 go build -o ./bin/oauth2-server ./cmd/oauth2-server/main.go

FROM scratch AS baseImage

WORKDIR /app

COPY --from=builder /app/bin/oauth2-server /app/oauth2-server
COPY --from=builder /app/private.key /app/private.key
COPY --from=builder /app/public.key /app/public.key
COPY --from=builder /app/.env /app/.env

EXPOSE 8080

CMD ["./oauth2-server"]
