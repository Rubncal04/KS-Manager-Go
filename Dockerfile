FROM golang:alpine

WORKDIR /go/src/ksmanager

RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum  ./

RUN go mod download


COPY . ./

EXPOSE 1323

CMD ["air", "-c", ".air.toml"]
