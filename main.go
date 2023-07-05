package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	cpuUsage := 0.2        // Desired CPU usage (e.g., 0.2 for 200m or 1.0 for 1000m)
	duration := 10 * time.Second // Duration for the CPU stress (e.g., 10 seconds)

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	numGoroutines := int(float64(numCPU) * cpuUsage)

	fmt.Printf("Starting CPU stress with %d goroutines...\n", numGoroutines)

	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(duration)
	close(done)

	fmt.Println("CPU stress completed.")
}
