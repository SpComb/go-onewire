FROM golang:1.9.4-stretch as go-build

RUN curl -L -o /tmp/dep-linux-amd64 https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && install -m 0755 /tmp/dep-linux-amd64 /usr/local/bin/dep

WORKDIR /go/src/github.com/SpComb/go-onewire

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . ./
RUN go install -v ./cmd/...



FROM debian:stretch

WORKDIR /opt/onewire
COPY --from=go-build /go/bin/w1-server /opt/onewire/bin/

CMD ["/opt/onewire/bin/w1-server", \
  "-verbose", \
  "-http-listen=:8286" \
]

EXPOSE 8286
