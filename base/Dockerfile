FROM golang:alpine AS builder

# Add source code
ADD ./ /go/src/gitlab.informatika.org/label-1-backend/base/

RUN cd /go/src/gitlab.informatika.org/label-1-backend/base && \       
    go build && \
    mv ./go-base /usr/bin/go-base

# Multi-Stage production build
FROM alpine

RUN apk add --update ca-certificates

# Retrieve the binary from the previous stage
COPY --from=builder /usr/bin/go-base /usr/local/bin/go-base

# Set the binary as the entrypoint of the container
CMD ["go-base", "serve"]