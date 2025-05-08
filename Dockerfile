FROM acgn0rg/tdlib:golang AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags "-s -w -extldflags '-static -fpic'" ./cmd/onest

FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /data

COPY --from=builder --chmod=755 /build/onest /usr/bin/onest

ENTRYPOINT [ "/usr/bin/onest" ]