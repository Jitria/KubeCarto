# SPDX-License-Identifier: Apache-2.0

### Builder

FROM golang:1.24-alpine3.21 as builder

RUN apk --no-cache update
RUN apk add --no-cache git clang llvm make gcc protobuf musl-dev

RUN mkdir /client
RUN mkdir /protobuf

WORKDIR /protobuf
COPY protobuf/ .

WORKDIR /client
COPY clients/log-client .

RUN go build -o log-client

### Make executable image

FROM alpine:3.17 as client

# RUN echo "@community http://dl-cdn.alpinelinux.org/alpine/edge/community" | tee -a /etc/apk/repositories
# RUN apk --no-cache update
# RUN apk add bash

COPY --from=builder /client/log-client /

CMD ["/log-client"]
