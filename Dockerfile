##### builder
FROM golang:1.14 as build
WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o app
      
##### result image
FROM alpine:3.10
COPY --from=build /go/src/app/app /app
COPY web web/
CMD /app