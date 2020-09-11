FROM golang:1.15-alpine as builder
WORKDIR /go/src/app
COPY src/ ./
RUN go build -o app main.go

FROM alpine:3.12
ENV GOTRACEBACK=single
COPY --from=builder /go/src/app .
CMD ["./app"]

