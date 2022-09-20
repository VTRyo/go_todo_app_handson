FROM golang:1.18.2-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app


# -----------

# container for deloy
FROM debian:bullseye-slim as deploy

RUN apt-get update
COPY --from=deploy-builder /app/app .

CMD ["./app"]


# ------------

# hot reload environment for use in local environment
FROM golang:1.18.2 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]