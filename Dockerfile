# Source: Go blog https://blog.golang.org/docker
# Source: https://medium.com/@monirz/golang-dependency-solution-with-go-module-and-docker-8967da6dd9f6
# Start Debian Image
FROM golang

# Setup Environment
ENV GO111MODULE=on

WORKDIR /app/server
# Copy go mod and go sum
COPY go.mod .
COPY go.sum .

# Copy rest of the file
COPY . .

# Install dependencies
RUN go mod download

# Build the app in Docker
WORKDIR /app/server/src
RUN go build -o server

EXPOSE 8080

CMD [ "./server" ]
