FROM docker.io/library/golang:1.23.7 AS builder

WORKDIR /go/src/github.com/aguidirh/go-rag-build
COPY . .
RUN NO_DOCKER=1 make build
RUN ls ./bin

FROM docker.io/library/golang:1.23.7
COPY --from=builder /go/src/github.com/aguidirh/go-rag-build/bin/go-rag-bot .
CMD go-rag-bot