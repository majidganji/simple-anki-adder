FROM golang:1.18.3-alpine3.16 as builder
RUN mkdir /app
WORKDIR /app
COPY . . 
RUN export GOPROXY=direct
RUN go mod download
RUN go build -o /app main.go

FROM scratch
COPY --from=builder /app /
CMD ["/app"]
