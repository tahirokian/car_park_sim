FROM golang:1.23.1-alpine

# Set working directory
RUN mkdir /vehicle_entry
WORKDIR /vehicle_entry

# Copy all files
COPY ./vehicle_entry /vehicle_entry

RUN go mod tidy

# Build the Go application
RUN go build -o vehicle_entry

# Command to run the Go service
CMD ["./vehicle_entry"]
