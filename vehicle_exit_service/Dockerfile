FROM golang:1.23.1-alpine

# Set working directory
RUN mkdir /vehicle_exit
WORKDIR /vehicle_exit

# Copy all files
COPY ./vehicle_exit /vehicle_exit

RUN go mod tidy

# Build the Go application
RUN go build -o vehicle_exit

# Command to run the Go service
CMD ["./vehicle_exit"]
