FROM golang:1.20.5-alpine3.18 as builder

WORKDIR /usr/src/app
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . .
RUN go build -o ./bin/app ./cmd

FROM alpine

WORKDIR /usr/src/app
COPY --from=builder ["/usr/src/app/bin/app", "/usr/src/app/"]
COPY --from=builder ["/usr/src/app/internal/db/migrations", "/usr/src/app/internal/db/migrations"]

EXPOSE 3000
CMD ["./app"]
