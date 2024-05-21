FROM golang:1.22.3-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./mouji .

FROM debian:bookworm-slim

COPY --from=builder /app/mouji /mouji

VOLUME /data
ENV DATA_FOLDER=/data

CMD ["/mouji"]
