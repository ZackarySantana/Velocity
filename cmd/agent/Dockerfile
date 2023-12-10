FROM golang:1.20 as vendor

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM vendor as build

COPY src src
COPY internal internal
COPY cmd/agent/main.go cmd/agent/main.go
RUN CGO_ENABLED=0 go build -o agent -ldflags '-extldflags "-static"' cmd/agent/main.go

FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY --from=build /app/agent /app/agent
RUN chmod +x /app/agent
ENV GIN_MODE=release
COPY velocity.yml ./
CMD ["/app/agent"]
