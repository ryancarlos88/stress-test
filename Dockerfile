# Use the official Go 1.19 base image
FROM golang:1.19-alpine as golang

# Set the working directory inside the container
WORKDIR /app
COPY . .

RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM scratch

# Copy the compiled Go application into the container
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=golang /app .
# Set the entry point for the container
ENTRYPOINT ["./app"]