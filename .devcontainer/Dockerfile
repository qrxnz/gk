FROM golang:1.24-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /usr/local/bin/gk .

FROM debian:bookworm-slim

RUN useradd --create-home --shell /bin/bash appuser

USER appuser
WORKDIR /home/appuser

COPY --from=build /usr/local/bin/gk /usr/local/bin/gk

CMD ["gk"]
