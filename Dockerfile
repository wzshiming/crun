FROM golang:alpine AS builder
WORKDIR /go/src/github.com/wzshiming/crun/
COPY . .
RUN go install github.com/wzshiming/crun/cmd/...

FROM wzshiming/upx AS upx
COPY --from=builder /go/bin/ /go/bin/
RUN upx /go/bin/*

FROM scratch
COPY --from=upx /go/bin/ /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/crun"]
