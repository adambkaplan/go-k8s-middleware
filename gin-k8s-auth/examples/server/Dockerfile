FROM registry.access.redhat.com/ubi9/go-toolset:latest

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY *.go ./
COPY examples examples
RUN go build -o ./main ./examples/server/main.go

EXPOSE 8081

CMD [ "/app/main" ]
