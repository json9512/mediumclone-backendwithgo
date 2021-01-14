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

# receive env variabes
ARG DBUsername
ARG DBHost
ARG DBPort
ARG DBName
ARG DBPassword

# Set env variables
ENV DB_USERNAME=${DBUsername}
ENV DB_HOST ${DBHost}
ENV DB_PORT ${DBPort}
ENV DB_NAME ${DBName}
ENV DB_PASSWORD ${DBPassword}

RUN go build -o server

EXPOSE 8080

CMD [ "./server" ]
