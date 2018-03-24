Linux onewire netlink golang library + REST API server.

# Support

Developed/tested on:

* Go 1.9
* Raspberry Pi 2 Model B, Raspbian stretch, Linux 4.9

# Libraries

### `github.com/SpComb/go-onewire/api`

API definitions for the `server`

[![](https://godoc.org/github.com/SpComb/go-onewire/api?status.svg)](http://godoc.org/github.com/SpComb/go-onewire/api)

### `github.com/SpComb/go-onewire/netlink/connector`

Linux `NETLINK_CONNECTOR` protocol support.

[![](https://godoc.org/github.com/SpComb/go-onewire/netlink/connector?status.svg)](http://godoc.org/github.com/SpComb/go-onewire/netlink/connector)

#### Issues
* https://github.com/SpComb/go-onewire/issues/1 netlink/connector requires forked github.com/mdlayher/netlink

### `github.com/SpComb/go-onewire/netlink/connector/w1`

Linux `w1` netlink protocol support.

[![](https://godoc.org/github.com/SpComb/go-onewire/netlink/connector/w1?status.svg)](http://godoc.org/github.com/SpComb/go-onewire/netlink/connector/w1)

### `github.com/SpComb/go-onewire/netlink/connector/w1/ds18b20`

DS18B20 specific support for the Linux `w1` netlink API.

[![](https://godoc.org/github.com/SpComb/go-onewire/netlink/connector/w1/ds18b20?status.svg)](http://godoc.org/github.com/SpComb/go-onewire/netlink/connector/w1/ds18b20)

# Server

## Usage
```
Usage of go/bin/w1-server:
  -debug
        Log debug
  -debug.server
        Log debug for server
  -debug.w1
        Log debug for w1
  -debug.web
        Log debug for web
  -http-listen string
        HTTP server listen: [HOST]:PORT (default ":8286")
  -http-static string
        HTTP sever /static path: PATH
  -quiet
        Do not log warnings
  -quiet.server
        Do not log warnings for server
  -quiet.w1
        Do not log warnings for w1
  -quiet.web
        Do not log warnings for web
  -verbose
        Log info
  -verbose.server
        Log info for server
  -verbose.w1
        Log info for w1
  -verbose.web
        Log info for web
```

## API

### `GET /api/`

```json
{
  "Sensors": [
    {
      "ID": "28-0315a4cfdbff",
      "Config": {
        "Bus": "onewire",
        "Type": "ds18b20",
        "Serial": "28-0315a4cfdbff"
      },
      "Status": {
        "At": "2018-03-24T19:44:12.488835929Z",
        "Temperature": 25
      }
    }
  ]
}
```
