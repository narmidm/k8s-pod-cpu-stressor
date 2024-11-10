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
)

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 20% CPU usage)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for CPU and memory stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU and memory stress indefinitely")
	memoryUsagePtr := flag.Float64("memory", 0.2, "Memory usage as a fraction (e.g., 0.2 for 20% memory usage)")
	flag.Parse()

	// Validate memory usage
	if *memoryUsagePtr < 0.0 || *memoryUsagePtr > 1.0 {
		fmt.Println("Error: --memory must be a value between 0.0 and 1.0 (representing 0% to 100% memory usage)")
		return
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Number of goroutines to use for stressing CPU
	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))
	if numGoroutines < 1 {
		numGoroutines = 1
	}

	// Log to indicate starting CPU and Memory stress
	fmt.Printf("Starting CPU stress with %d goroutines targeting %.2f%% CPU usage and targeting %.2f%% memory usage...\n", numGoroutines, *cpuUsagePtr*100, *memoryUsagePtr*100)

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	var stopFlag int32

	// CPU stress logic
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
					_ = rand.Float64() * rand.Float64()
				}

				// Idle for the rest of the interval
				time.Sleep(idleDuration)
			}
		}()
	}

	// Memory stress logic
	go func() {
		for {
			if atomic.LoadInt32(&stopFlag) == 1 {
				return
			}

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			memUsage := int(float64(memStats.Sys) * (*memoryUsagePtr)) // Use Sys to get the total allocated from OS
			if memUsage < 1 {
				memUsage = 1
			}

			memStress := make([]byte, memUsage)
			for i := range memStress {
				memStress[i] = byte(rand.Intn(256))
			}

			// Simulate memory usage
			_ = memStress[rand.Intn(len(memStress))]

			// Sleep for a short duration to simulate memory stress
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Wait for either the quit signal or the duration to elapse
	go func() {
		<-quit
		fmt.Println("\nTermination signal received. Stopping CPU and memory stress...")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		fmt.Println("\nCPU and memory stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
		// Keep the process running to prevent pod from restarting
		select {}
	}

	// Run stress indefinitely
	fmt.Println("CPU and memory stress will run indefinitely. Press Ctrl+C to stop.")
	<-done
}
