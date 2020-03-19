FROM golang:1.13.5-alpine as builder
LABEL 9to6="ktk0011@gmail.com"

RUN apk add --update git make

ADD . /aws-ec2-price
WORKDIR /aws-ec2-price

ENV CGO_ENABLED 0
RUN go install -mod vendor ./cmd/main.go

FROM alpine:3
COPY --from=builder /go/bin/main /bin/aws-ec2-price

EXPOSE 8080

ENTRYPOINT ["aws-ec2-price"]
