FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build

# Copy and download dependency using go mod
ADD . /go/src/myapp
WORKDIR /go/src/myapp

RUN go mod download

# Copy the code into the container

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go/src/myapp/main .

# Export necessary port
EXPOSE 6000

# Command to run when starting the container
CMD ["/dist/main"]
