package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "stress-test",
		Short: "Perform stress tests",
		Run:   runStressTest,
	}

	rootCmd.Flags().StringP("url", "u", "", "Duration of the stress test in seconds")
	rootCmd.Flags().IntP("requests", "r", 10, "Number of requests to perform")
	rootCmd.Flags().IntP("concurrency", "c", 10, "Number of concurrent requests")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func runStressTest(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("url")
	numRequests, _ := cmd.Flags().GetInt("requests")
	concurrency, _ := cmd.Flags().GetInt("concurrency")

	fmt.Printf("Starting stress test with %d requests and %d concurrent\n", numRequests, concurrency)

	performStressTest(url, numRequests, concurrency)

	fmt.Println("Stress test completed")
}

func performStressTest(url string, numRequests int, concurrency int) {
	timer := time.Now()

	results := make(chan int)
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < numRequests/concurrency; j++ {
				success := performRequest(url, &wg)
				results <- success
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make(map[int]int)
	for result := range results {
		_, ok := out[result]
		if !ok {
			out[result] = 1
		} else {
			out[result]++

		}
	}
	fmt.Printf("\nStress test took us %.2f seconds with the following results\n\n", time.Since(timer).Seconds())
	fmt.Printf("Number of requests with status code 200 : %d\n\n", out[200])
	fmt.Printf("Here are all the %d different status codes we got for all %d requets:\n\n", len(out), numRequests)
	fmt.Printf("Status Code | Occurences\n")
	for k, v := range out {
		fmt.Printf("    %d     |     %d\n", k, v)
	}
	fmt.Println()
}

func performRequest(url string, wg *sync.WaitGroup) int {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
