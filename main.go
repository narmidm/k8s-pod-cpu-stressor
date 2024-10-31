package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 20% CPU usage)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU stress indefinitely")
	flag.Parse()

	// Validate input parameters
	if *cpuUsagePtr <= 0 || *cpuUsagePtr > 1 {
		log.Fatalf("Invalid CPU usage: %f. It must be between 0 and 1.", *cpuUsagePtr)
	}

	if *durationPtr <= 0 {
		log.Fatalf("Invalid duration: %s. It must be greater than 0.", *durationPtr)
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Number of goroutines to use for stressing CPU
	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))
	if numGoroutines < 1 {
		numGoroutines = 1
	}

	log.Infof("Starting CPU stress with %d goroutines targeting %.2f CPU usage...", numGoroutines, *cpuUsagePtr)

	done := make(chan struct{})

	// Capture termination signals
	quit := make(chan os.Signal, 1)
	if err := signal.Notify(quit, os.Interrupt, os.Kill); err != nil {
		log.Fatalf("Failed to set up signal notification: %v", err)
	}

	var stopFlag int32

	// Improved workload generation
	for i := 0; i < numGoroutines; i++ {
		go func() {
			workDuration := time.Duration(*cpuUsagePtr*1000) * time.Microsecond
			idleDuration := time.Duration((1-*cpuUsagePtr)*1000) * time.Microsecond

			for {
				if atomic.LoadInt32(&stopFlag) == 1 {
					return
				}

				// Busy loop for the specified work duration
				endWork := time.Now().Add(workDuration)
				for time.Now().Before(endWork) {
					// Perform a small computation to keep the CPU active
					_ = rand.Float64() * rand.Float64()
				}

				// Idle for the rest of the interval
				time.Sleep(idleDuration)
			}
		}()
	}

	go func() {
		// Wait for termination signal
		<-quit
		log.Println("\nTermination signal received. Stopping CPU stress...")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		log.Println("\nCPU stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
		// Keep the process running to prevent the pod from restarting
		select {}
	}

	// Run stress indefinitely
	log.Println("CPU stress will run indefinitely. Press Ctrl+C to stop.")
	<-done
}
