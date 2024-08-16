# builder stage
FROM golang:1.22.4-alpine as builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /opt

COPY go.mod go.sum pkg internal ./

RUN go mod download

RUN go build -o main cmd/tm-player/main.go


# runtime stage
FROM scratch as runtime

WORKDIR /opt

COPY --from=builder /opt/main /main

CMD ["/main"]
