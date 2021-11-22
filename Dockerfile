# Start from golang base image
FROM golang:1.16 as builder

WORKDIR /go/src/github.com/ONSdigital/ssdc-rm-eq-launcher

COPY . .

# Download dependencies
RUN go get

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod mod -o /go/bin/ssdc-rm-eq-launcher .

######## Start a new stage from scratch #######
FROM alpine:latest  

# Copy the Pre-built binary file and entry point from the previous stage
COPY --from=builder /go/bin/ssdc-rm-eq-launcher .
COPY docker-entrypoint.sh .
COPY static/ /static/
COPY jwt-test-keys /jwt-test-keys/

EXPOSE 8000

ENTRYPOINT ["sh", "/docker-entrypoint.sh"]
