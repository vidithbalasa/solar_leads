package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	// Open the file containing addresses
	inputFile, err := os.Open("addresses.txt")
	if err != nil {
		fmt.Println("Error opening addresses file:", err)
		return
	}
	defer inputFile.Close()

	// Create the output file to store results
	outputFile, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}
	defer outputFile.Close()

	// Channel to pass addresses to goroutines
	addressChannel := make(chan string)
	// Channel to collect results
	resultChannel := make(chan string)

	// Number of concurrent workers
	const numWorkers = 10
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go worker(addressChannel, resultChannel, &wg)
	}

	// Goroutine to collect results and write to file
	go func() {
		for result := range resultChannel {
			outputFile.WriteString(result)
		}
	}()

	// Read addresses from file and send them to workers
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		addressChannel <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input file:", err)
	}

	// Close the address channel and wait for workers to finish
	close(addressChannel)
	wg.Wait()
	close(resultChannel)
}

// worker processes addresses from the addressChannel, makes HTTP requests, and sends results to the resultChannel
func worker(addressChannel <-chan string, resultChannel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for address := range addressChannel {
		cleanAddress := cleanAddress(address)
		requestURL := fmt.Sprintf("http://localhost:8080/getSolarData?address=%s", url.QueryEscape(cleanAddress))
		response, err := http.Get(requestURL)
		if err != nil {
			resultChannel <- fmt.Sprintf("%s\nError: %v\n", cleanAddress, err)
			continue
		}
		responseBody, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			resultChannel <- fmt.Sprintf("%s\nError reading response: %v\n", cleanAddress, err)
			continue
		}
		resultChannel <- fmt.Sprintf("%s\n%s\n", cleanAddress, string(responseBody))
	}
}

// cleanAddress removes spaces and other unnecessary characters from the address
func cleanAddress(address string) string {
	return strings.Join(strings.Fields(address), "")
}

