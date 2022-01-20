FROM golang:1.17.6-buster

WORKDIR /app

# Copy the go source:
COPY . .

# Copy the go build artifact:
COPY bin/pa-go-starter-amd64-linux server

# Go get dependencies.  
RUN go get -d ./...
RUN go install -v ./...

CMD ["/app/server"]
