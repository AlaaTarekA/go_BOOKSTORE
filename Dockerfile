FROM golang:1.17.1

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/golang-sample-app-master


# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8000 to the outside world
EXPOSE 8000

# Run the executable
CMD ["go-sample-app"]

