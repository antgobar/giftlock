FROM golang:1.24.2-alpine3.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o /bin/giftlock cmd/web/main.go

FROM alpine:latest AS final
RUN apk update && apk add --no-cache ca-certificates
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

WORKDIR /app
USER appuser

COPY --from=build /bin/giftlock .
COPY --from=build /app/config ./config
COPY --from=build /app/assets ./assets
COPY --from=build /app/sql ./sql
EXPOSE 8080
ENTRYPOINT ["sh", "-c", "./giftlock"]