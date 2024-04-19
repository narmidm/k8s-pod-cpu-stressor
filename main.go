package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 200m)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU stress indefinitely")
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))

	fmt.Printf("Starting CPU stress with %d goroutines...\n", numGoroutines)

	done := make(chan struct{})

	// Capture termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

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

	go func() {
		// Wait for termination signal
		<-quit
		fmt.Println("Termination signal received. Stopping CPU stress...")
		close(done)
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		fmt.Println("CPU stress completed.")
		os.Exit(0)
	}

	// Run stress indefinitely
	fmt.Println("CPU stress will run indefinitely. Press Ctrl+C to stop.")
	select {}
}
