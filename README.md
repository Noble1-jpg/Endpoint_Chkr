# Endpoint Checker

This tool checks the availability of a list of HTTP endpoints provided in a YAML configuration file. It evaluates each endpoint every 15 seconds and logs availability statistics by domain.

## Table of Contents

- [What It Does](#what-it-does)
- [Installation](#installation)
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
example.com has 100% availability
```

This output appears every 15 seconds.

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

### Request body handling
- **Issue:** Sent full struct as JSON
- **Original Behavior:** Entire endpoint struct was marshaled and sent in request body
- **Fix Implemented:** Now only the `body` field is used for request payload

### No timeout
- **Issue:** Could hang on slow endpoints
- **Original Behavior:** HTTP requests had no time limit
- **Fix Implemented:** Timeout of 500ms added to the HTTP client

### Incomplete availability check
- **Issue:** Only status code was checked
- **Original Behavior:** Ignored response time in availability logic
- **Fix Implemented:** Now also checks response duration to ensure it's within 500ms

### Crashes from nil map entries
- **Issue:** Map entries accessed without initialization
- **Original Behavior:** Could panic if domain key didn't exist
- **Fix Implemented:** Map values are initialized safely using nil checks

### Domain detection
- **Issue:** Included ports in domain name
- **Original Behavior:** Used full host (including port)
- **Fix Implemented:** Now extracts only hostname, ignoring port

### Sequential execution
- **Issue:** Slow with many endpoints
- **Original Behavior:** Endpoints checked one after another
- **Fix Implemented:** Uses goroutines and WaitGroup for concurrent execution

### Unsafe shared state
- **Issue:** Risk of data races on stats map
- **Original Behavior:** Map was updated from multiple goroutines without protection
- **Fix Implemented:** Mutex added to guard map access

### Decimal percentages
- **Issue:** Showed float availability values
- **Original Behavior:** Percentages were printed with decimals
- **Fix Implemented:** Availability now rounded down to whole numbers

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



