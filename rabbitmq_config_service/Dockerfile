FROM golang:1.23.1-alpine

# Set working directory
RUN mkdir /rabbitmq_configurator
WORKDIR /rabbitmq_configurator

# Copy all files
COPY ./rabbitmq_configurator /rabbitmq_configurator

RUN go mod tidy

# Build the Go application
RUN go build -o rabbitmq_configurator

# Command to run the Go service
CMD ["./rabbitmq_configurator"]
