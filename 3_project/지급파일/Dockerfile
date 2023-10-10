FROM public.ecr.aws/docker/library/golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

COPY main.go ./

RUN go build -o main .

FROM public.ecr.aws/docker/library/alpine

WORKDIR /app

COPY /build/main .

CMD ["./main"]