FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build

# Copy the code into the container
ADD . /go/src/myapp
WORKDIR /go/src/myapp

# Copy and download dependency using go mod
RUN go mod download

# Build the application
RUN go build -o main.exe main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go/src/myapp/main.exe .

# Export necessary port
EXPOSE 6000

# Command to run when starting the container
CMD ["/dist/main.exe"]
