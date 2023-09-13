#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .
COPY go.sum .

RUN go mod download

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o main /app/cmd/peoples_graph/main.go

#####################################
#   STEP 2 build a small image      #
#####################################
FROM scratch

# Copy our static executable.
COPY --from=builder /app/main /app/main
COPY --from=builder /app/deploy/peoples_graph/.env /app/deploy/peoples_graph/.env

# Run the hello binary.
ENTRYPOINT ["/app/main"]