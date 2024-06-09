FROM golang:1.22.2

# Set the working directory in the container
WORKDIR /src

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Read version from build-args
ARG VERSION
RUN echo "Build version $VERSION"

# Set the working directory to the cmd directory
WORKDIR /src/cmd

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /api .

# Expose port 3000
EXPOSE 3000

# Command to run the application
CMD ["/api"]
