# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
#FROM golang:1.12-alpine
FROM golang:1.12.0-alpine3.9 as builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh

# Add Maintainer Info
#LABEL maintainer="Ricky Lin <rick3325@gmail.com>"

# Set the Current Working Directory inside the container
ADD . /app
WORKDIR /app


# Copy go mod and sum files
#COPY go.mod go.sum ./


# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN apk add git
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
#COPY . .

#RUN cd /app && go build
RUN cd /app
RUN go mod init httpapitest
#RUN go get -u github.com/gorilla/mux

RUN go build

# Build the Go app
#RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8082

# Run the executable
#CMD ["./main"]
ENTRYPOINT ./httpapitest