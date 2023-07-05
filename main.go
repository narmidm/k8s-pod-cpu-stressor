package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 200m)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))

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

	time.Sleep(*durationPtr)
	close(done)

	fmt.Println("CPU stress completed.")
}
