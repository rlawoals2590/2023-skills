FROM public.ecr.aws/docker/library/golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \ 
    GOARCH=amd64

WORKDIR /build

COPY main.go ./

COPY templates/ ./templates/

RUN go mod init example.com/m

RUN go mod tidy

RUN go build -o main .

EXPOSE 8080

FROM public.ecr.aws/docker/library/alpine

WORKDIR /app

COPY --from=builder /build/ .

CMD ["./main"]