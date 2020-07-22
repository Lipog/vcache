FROM alpine:latest
COPY . /usr/bin/main
RUN chmod 544 /usr/bin/main
WORKDIR /usr/bin/
ENTRYPOINT ["/bin/sh","-c","main"]

