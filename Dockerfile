FROM golang:1.21 AS build
WORKDIR /go/src/ewallet/
COPY . /go/src/ewallet/
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -installsuffix cgo -o /go/src/ewallet/build/ewallet /go/src/ewallet/cmd/app/main.go

FROM busybox AS runtime
WORKDIR /app
COPY --from=build /go/src/ewallet/build/ewallet /app/
EXPOSE 8080/tcp
ENTRYPOINT ["./ewallet"]
