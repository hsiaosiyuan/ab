## What's it
a golang toy similar ab(apache benchmarking tool)

## Usage

    go run ab.go -h

    -C=&[]: Cookie. add cookie to the request (repeatable)
    -E="": Url of the endpoint for testing
    -H=&[]: Custom header, append extra headers to the request (repeatable)
    -T="application/x-www-form-urlencoded": Content-type header to use for POST/PUT data, eg. application/x-www-form-urlencoded. default: text/plain.
    -c=1: Number of multiple requests to perform at a time
    -e="": Parameter name will be appended in query string to escape cache. eg: '_ec=timestamp'
    -n=1: Number of requests to perform for the benchmarking session
    -p=: File containing data to POST, remember to also set -T
    -t=3600: Time limit. maximum number of seconds to spend for benchmarking

## Example

**normal**

    go run ab.go -H="test: test" -C="test=test" -H="test1:test1" -C="test1=test1" -n=10 -c=2 -E='http://example.com'

output

    Endpoint to test: http://example.com

    Complete  2  requests
    Complete  4  requests
    Complete  6  requests
    Complete  8  requests
    Complete  10  requests

    All requests: 10
    Time taken: 5.286  [second]
    Succeed requests: 10
    Failed requests: 0
    Non2xx requests: 0
    Body sent: 0  [bytes]
    HTML transferred: 12700  [bytes]
    Request per second: 1.89  [#/sec] (mean)

**escape cache**

    go run ab.go -e='_ec' -H="test: test" -C="test=test" -H="test1:test1" -C="test1=test1" -n=10 -c=2 -E='http://example.com'

output

    Endpoint to test: http://example.com  (will escape cache)

    Complete  2  requests
    Complete  4  requests
    Complete  6  requests
    Complete  8  requests
    Complete  10  requests

    All requests: 10
    Time taken: 3.644  [second]
    Succeed requests: 10
    Failed requests: 0
    Non2xx requests: 0
    Body sent: 0  [bytes]
    HTML transferred: 12700  [bytes]
    Request per second: 2.74  [#/sec] (mean)
