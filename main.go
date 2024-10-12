package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 20% CPU usage)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU stress indefinitely")
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Number of goroutines to use for stressing CPU
	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))
	if numGoroutines < 1 {
		numGoroutines = 1
	}

	fmt.Printf("Starting CPU stress with %d goroutines targeting %.2f CPU usage...\n", numGoroutines, *cpuUsagePtr)

	done := make(chan struct{})

	// Capture termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	var stopFlag int32

	for i := 0; i < numGoroutines; i++ {
		go func() {
			workDuration := time.Duration(*cpuUsagePtr*100) * time.Millisecond
			idleDuration := time.Duration((1-*cpuUsagePtr)*100) * time.Millisecond

			for {
				if atomic.LoadInt32(&stopFlag) == 1 {
					return
				}

				// Busy loop for the specified work duration
				endWork := time.Now().Add(workDuration)
				for time.Now().Before(endWork) {
				}

				// Idle for the rest of the interval
				time.Sleep(idleDuration)
			}
		}()
	}

	go func() {
		// Wait for termination signal
		<-quit
		fmt.Println("\nTermination signal received. Stopping CPU stress...")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		fmt.Println("\nCPU stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
		os.Exit(0)
	}

	// Run stress indefinitely
	fmt.Println("CPU stress will run indefinitely. Press Ctrl+C to stop.")
	<-done
}
