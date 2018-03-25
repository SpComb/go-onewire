FROM golang:1.9.4-stretch as go-build

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/SpComb/go-onewire

COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . ./
RUN go install -v ./cmd/...

CMD ["/go/bin/w1-server", \
  "-verbose", \
  "-http-listen=:8286" \
]

EXPOSE 8286
