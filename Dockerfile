# Build stage
FROM golang:1.21.0-alpine3.17 AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download
# Copia el archivo .env
COPY .env .  
COPY *go .
RUN go build -o app 

# Define las variables de entorno para Firebase
ENV GOOGLE_APPLICATION_CREDENTIALS="/app/.env"

EXPOSE 8080
ENTRYPOINT [ "./app" ]

