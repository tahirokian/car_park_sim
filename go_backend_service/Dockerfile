FROM golang:1.23.2-alpine

# Set working directory
RUN mkdir /go_backend
WORKDIR /go_backend

# Copy all files
COPY ./go_backend /go_backend

RUN go mod tidy

# Build the Go application
RUN go build -o go_backend

# Command to run the Go service
CMD ["./go_backend"]
