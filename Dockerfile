FROM golang:1.13.4-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy files
COPY *.go ./
COPY go.mod .
COPY go.sum .
COPY .env .

# Necessary to avoid needing gcc
ENV CGO_ENABLED=0

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main .

# Run the executable
CMD ["./main"]
