FROM golang:latest

WORKDIR /blackbox

# Copy only go.mod file to check if dependencies have changed
COPY go.mod .

# Download modules while utilizing cache
RUN go mod download

COPY . .

# Build the Go app with optimizations enabled
RUN CGO_ENABLED=0 GOOS=linux go build -o main -a -installsuffix cgo . && \
    strip --strip-unneeded main

# Use scratch as a base image to reduce the final image size
FROM scratch

COPY --from=0 /blackbox/main .

# Expose port 8000 to the outside world
EXPOSE 8000

CMD ["/main"]
