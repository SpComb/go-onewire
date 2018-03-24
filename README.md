
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
