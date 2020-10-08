# STEP 1 build binary
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY . /app/

# Fetch dependencies.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# STEP 2 build a small image
FROM scratch

# Copy static executable and certificates
COPY --from=builder /app/main /app/main
COPY /www2 /app/www2
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY extras/items.csv.gz /app/

WORKDIR /app
# Run the binary.
ENTRYPOINT ["/app/main", "--csv", "items.csv.gz"]
