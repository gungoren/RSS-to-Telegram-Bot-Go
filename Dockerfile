FROM golang:1.12-alpine AS build
WORKDIR /src/app
COPY . .
RUN go build -o /app


FROM alpine
COPY --from=build /app /app

# Define environment variable
ENV TOKEN X
ENV CHATID X
ENV DELAY 60

ENTRYPOINT [ "/app" ]