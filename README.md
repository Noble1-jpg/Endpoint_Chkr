# Endpoint Checker

This tool checks the availability of a list of HTTP endpoints provided in a YAML configuration file. It evaluates each endpoint every 15 seconds and logs availability statistics by domain.

## Table of Contents

- [What does it do?](#what-does-it-do)
- [Installation Steps](#installation-steps)
- [Running the Program](#running-the-program)
- [Configuration File Format](#configuration-file-format)
- [Identified Issues and Fixes](#identified-issues-and-fixes)
- [What I Learned and Enjoyed](#what-i-learned-and-enjoyed)
- [Production-Readiness](#production-readiness)
- [License](#license)
- [Author](#author)

## What It Does

- Reads endpoints from a YAML file.
- Sends requests based on provided HTTP method, headers, and body.
- Determines if an endpoint is available:
  - If it responds with a status code between 200 and 299.
  - If it responds within 500 milliseconds.
- Tracks and logs cumulative availability by domain (ignoring port numbers).
- Executes checks in parallel for efficiency.

## Installation

1. Install Go if you haven't already: https://golang.org/doc/install
2. Clone this repository:

```
git clone https://github.com/Noble1-jpg/Endpoint_Chkr.git
cd Endpoint_Chkr
```

3. Build the application:

```
go build -o endpoint-checker main.go
```

## Running the Program

To run the checker, use:

```
./endpoint-checker config.yaml
```

You will see output like:

```
Nabilah.com has 100% availability
```

This output will appear every 15 seconds.

## Configuration File Format

The configuration file should be in YAML format. Example:

```yaml
- name: check homepage
  url: https://example.com/

- name: post example
  url: https://example.com/api
  method: POST
  headers:
    content-type: application/json
  body: '{"hello":"world"}'
```

## Identified Issues and Fixes

- Request body was wrong: It sent the whole endpoint struct as JSON. Now it just sends the actual body field.
- No timeout on requests: If a server was slow, the program could freeze. Now it gives up after 500 milliseconds.
- Didn’t check response time: Only status code was checked. Now it checks that and how fast the response was.
- Crashes with missing map entries: It tried to update the stats map without setting it up first. Now it safely creates new entries.
- Included port in domain name: It treated "example.com:443" and "example.com" as different. Now it only uses the hostname.
- Ran endpoints one after another: This was too slow. Now it runs them all at the same time using goroutines.
- Shared data could cause bugs: Multiple goroutines could access shared stats at the same time. A mutex now keeps this safe.
- Availability was a decimal: It showed numbers like 98.3%. Now it rounds down to whole numbers like 98%.


## What I Learned and Enjoyed

- I learned how to handle concurrency safely in Go using `sync.Mutex` and `sync.WaitGroup`.
- I gained more experience parsing and validating YAML files.
- I better understand HTTP client configuration, including timeouts and headers.
- I enjoyed improving performance by adding parallelism.
- I also enjoyed debugging and making the program resilient to errors.

## Production-Readiness

- Safe concurrency with mutex locking
- Graceful handling of request errors
- Consistent logging every 15 seconds
- Proper parsing and validation of configuration input

## License

MIT License

## Author

Nabilah Aralepo - (https://github.com/Noble1-jpg)


