#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk --no-cache add \
          alpine-sdk \
          librdkafka-dev \
          pkgconf && \
        rm -rf /var/cache/apk/*


WORKDIR /app

COPY . .
COPY go.sum .

RUN go mod download

# Build the binary.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-extldflags -static" -tags musl \
     -o main /app/cmd/peoples_kafka/main.go

#####################################
#   STEP 2 build a small image      #
#####################################
FROM scratch

# Copy our static executable.
COPY --from=builder /app/main /app/main
COPY --from=builder /app/deploy/peoples_kafka/.env /app/deploy/peoples_kafka/.env

# Run the hello binary.
ENTRYPOINT ["/app/main"]