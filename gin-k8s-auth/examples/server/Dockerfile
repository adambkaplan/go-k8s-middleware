FROM registry.redhat.io/ubi9/go-toolset:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY examples examples

RUN go build -o ./server ./examples/server/main.go

CMD [ "/app/server" ]