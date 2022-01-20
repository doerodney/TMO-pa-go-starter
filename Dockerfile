# Use buster so /bin/bash is avsailable:
FROM golang:1.17.6-buster

# Make the /app directory:
RUN mkdir /app

# Put all the local src into /app:
ADD . /app

# Specify /app as the working directory:
WORKDIR /app

# Build locally in the container - no cross-compile necessary:
# Built executable is "server"
RUN go build -o server .

# Run /app/server:
CMD ["/app/server"]
