package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type DomainStats struct {
	Success int
	Total   int
}

var (
	statsMu sync.Mutex
	stats   = make(map[string]*DomainStats)
)

func checkHealth(endpoint Endpoint) {
	client := &http.Client{Timeout: 500 * time.Millisecond}

	var reqBody *bytes.Reader
	if endpoint.Body != "" {
		reqBody = bytes.NewReader([]byte(endpoint.Body))
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(endpoint.Method, endpoint.URL, reqBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	for k, v := range endpoint.Headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	domain := extractDomain(endpoint.URL)

	statsMu.Lock()
	defer statsMu.Unlock()
	if stats[domain] == nil {
		stats[domain] = &DomainStats{}
	}
	stats[domain].Total++

	if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 && duration <= 500*time.Millisecond {
		stats[domain].Success++
	}
}

func extractDomain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	host := u.Hostname() // excludes port
	return host
}

func monitorEndpoints(endpoints []Endpoint) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		var wg sync.WaitGroup
		for _, endpoint := range endpoints {
			wg.Add(1)
			go func(ep Endpoint) {
				defer wg.Done()
				checkHealth(ep)
			}(endpoint)
		}
		wg.Wait()
		logResults()
		<-ticker.C
	}
}

func logResults() {
	statsMu.Lock()
	defer statsMu.Unlock()
	for domain, stat := range stats {
		availability := 0
		if stat.Total > 0 {
			availability = int(float64(stat.Success) / float64(stat.Total) * 100)
		}
		fmt.Printf("%s has %d%% availability\n", domain, availability)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config_file.yaml>")
	}

	filePath := os.Args[1]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var endpoints []Endpoint
	if err := yaml.Unmarshal(data, &endpoints); err != nil {
		log.Fatal("Error parsing YAML:", err)
	}

	monitorEndpoints(endpoints)
}
