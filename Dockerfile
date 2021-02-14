FROM golang:1.14 as builder

# Set Environment Variables
ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o main .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

# Define environment variable
ENV TOKEN X
ENV CHATID X
ENV DELAY 60

ENTRYPOINT [ "./main" ]