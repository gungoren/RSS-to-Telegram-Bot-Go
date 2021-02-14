FROM golang:1.14 as builder

# Set Environment Variables
ENV HOME /app
ENV CGO_ENABLED 1
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -a -ldflags "-linkmode external -extldflags -static" -o main .

FROM alpine:edge

RUN apk --update upgrade
RUN apk add sqlite
# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# removing apk cache
RUN rm -rf /var/cache/apk/*

COPY --from=builder /app/main .

# Define environment variable
ENV TOKEN X
ENV CHATID X
ENV DELAY 180

ENTRYPOINT [ "./main" ]