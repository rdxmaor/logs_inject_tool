FROM golang:1.12 as build

WORKDIR /build

COPY . .
RUN go build -o /build/bin/logs_inject_tool /build/cmd/main.go

FROM ubuntu:18.04 as final

WORKDIR /app
COPY --from=build /build/bin/ .

CMD ./logs_inject_tool