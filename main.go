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
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU stress indefinitely")
	memoryUsagePtr := flag.Float64("memory", 0.2, "Memory usage as a fraction (e.g., 0.2 for 20% memory usage)")
	memDurationPtr := flag.Duration("memduration", 10*time.Second, "Duration for the memory stress (e.g., 10s)")
	memForeverPtr := flag.Bool("memforever", false, "Run memory stress indefinitely")
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

	// Memory stress logic
	go func() {
		memUsage := int(float64(runtime.MemStats.Sys) * (*memoryUsagePtr))
		memStress := make([]byte, memUsage)
		for i := range memStress {
			memStress[i] = byte(rand.Intn(256))
		}

		for {
			if atomic.LoadInt32(&stopFlag) == 1 {
				return
			}

			// Simulate memory usage
			_ = memStress[rand.Intn(len(memStress))]

			// Sleep for a short duration to simulate memory stress
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		// Wait for termination signal
		<-quit
		fmt.Println("\nTermination signal received. Stopping CPU and memory stress...")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		fmt.Println("\nCPU stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
		// Keep the process running to prevent the pod from restarting
		select {}
	}

	if !*memForeverPtr {
		time.Sleep(*memDurationPtr)
		fmt.Println("\nMemory stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		close(done)
		// Keep the process running to prevent the pod from restarting
		select {}
	}

	// Run stress indefinitely
	fmt.Println("CPU and memory stress will run indefinitely. Press Ctrl+C to stop.")
	<-done
}
