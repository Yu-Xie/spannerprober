# ------ builder ------
FROM golang:1.10 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/code.uber.internal/marketplace/spannerprober

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix nocgo -o /app .


# ------ runner ------
FROM alpine
RUN apk add bash
RUN apk add --no-cache ca-certificates apache2-utils
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache tcpdump

# copy binary from builder
COPY --from=builder /app ./

ENTRYPOINT ["./app"]
