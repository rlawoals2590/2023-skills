FROM public.ecr.aws/docker/library/golang:1.16.15

WORKDIR /app

RUN go mod init example.com/m
RUN go mod download

COPY product ./

# Build
RUN chmod 777 ./product


RUN adduser product

USER product
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["./product"]
