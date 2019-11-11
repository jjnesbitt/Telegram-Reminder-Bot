FROM golang:1.13.4-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy files
COPY *.go ./
COPY .env .

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Install dependencies
RUN go get github.com/joho/godotenv
RUN go get gopkg.in/tucnak/telebot.v2

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 9000

# Run the executable
CMD ["./main"]
