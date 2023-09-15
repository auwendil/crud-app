FROM golang:1.20

WORKDIR /var/crud-app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /opt ./...
RUN mv /opt/api /opt/crud-app

CMD ["/opt/crud-app"]