package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	version = "1.0.0"
)

// CPUUsageMonitor tracks CPU usage and provides feedback
type CPUUsageMonitor struct {
	targetUsage    float64
	currentUsage   float64
	adjustmentLock sync.Mutex
	scaleFactor    float64 // Adjustment multiplier for workload
}

// NewCPUUsageMonitor creates a new CPU usage monitor
func NewCPUUsageMonitor(targetUsage float64) *CPUUsageMonitor {
	return &CPUUsageMonitor{
		targetUsage:  targetUsage,
		scaleFactor:  1.0, // Start with no adjustment
		currentUsage: 0,
	}
}

// AdjustWorkload returns the current scale factor for the workload
func (m *CPUUsageMonitor) AdjustWorkload() float64 {
	m.adjustmentLock.Lock()
	defer m.adjustmentLock.Unlock()
	return m.scaleFactor
}

// UpdateUsage updates the current CPU usage and adjusts the scale factor
func (m *CPUUsageMonitor) UpdateUsage(actualUsage float64) {
	m.adjustmentLock.Lock()
	defer m.adjustmentLock.Unlock()

	m.currentUsage = actualUsage

	// Simple proportional control - adjust based on ratio of target to actual
	if actualUsage > 0.01 { // Avoid division by very small numbers
		// If we're using too much CPU, decrease the scale factor
		// If we're using too little, increase it
		adjustment := m.targetUsage / actualUsage

		// Limit adjustment rate to avoid oscillation
		if adjustment > 2.0 {
			adjustment = 2.0
		} else if adjustment < 0.5 {
			adjustment = 0.5
		}

		// Gradually adjust the scale factor with stronger weight for recent measurements
		m.scaleFactor = m.scaleFactor*0.5 + adjustment*0.5
	}

	fmt.Printf("CPU Usage: %.2f%% (target: %.2f%%), adjustment factor: %.2f\n",
		actualUsage*100, m.targetUsage*100, m.scaleFactor)
}

// getCPUUsage returns a relative measure of CPU performance
func getCPUUsage() float64 {
	// Get the start time
	start := time.Now()
	var iterations uint64

	// Run some work to measure how many ops/sec we can do
	for i := 0; i < 1000000; i++ {
		iterations++
		// Do some meaningless work
		_ = math.Sqrt(rand.Float64())
	}
	elapsed := time.Since(start)

	// Calculate the CPU usage based on how much work we accomplished
	return float64(iterations) / float64(elapsed.Nanoseconds())
}

func printUsage() {
	fmt.Printf("k8s-pod-cpu-stressor %s\n", version)
	fmt.Println("A tool for simulating CPU load in Kubernetes pods")
	fmt.Println("\nUsage:")
	fmt.Println("  cpu-stress [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -cpu=0.2       CPU usage as a fraction (e.g., 0.2 for 20% CPU usage)")
	fmt.Println("  -duration=10s   Duration for the CPU stress (e.g., 10s, 5m, 1h)")
	fmt.Println("  -forever        Run CPU stress indefinitely")
	fmt.Println("  -version        Show version information")
	fmt.Println("  -help           Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  cpu-stress -cpu=0.5 -duration=30s    # Use 50% CPU for 30 seconds")
	fmt.Println("  cpu-stress -cpu=0.8 -forever         # Use 80% CPU indefinitely")
}

func main() {
	// Parse command-line flags
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 20% CPU usage)")
	durationPtr := flag.Duration("duration", 10*time.Second, "Duration for the CPU stress (e.g., 10s)")
	runForeverPtr := flag.Bool("forever", false, "Run CPU stress indefinitely")
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help message")

	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Printf("k8s-pod-cpu-stressor version %s\n", version)
		os.Exit(0)
	}

	// Show help if requested
	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	// Validate CPU usage
	if *cpuUsagePtr <= 0 || *cpuUsagePtr > 1.0 {
		fmt.Printf("Error: CPU usage must be between 0 and 1.0, got %.2f\n", *cpuUsagePtr)
		os.Exit(1)
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// Get baseline CPU measurement for calibration
	baselineCPU := getCPUUsage()
	fmt.Printf("Baseline CPU measurement: %.6f ops/ns\n", baselineCPU)

	// Create CPU usage monitor with the target usage
	monitor := NewCPUUsageMonitor(*cpuUsagePtr)

	// Number of goroutines to use for stressing CPU
	numGoroutines := int(float64(numCPU)*(*cpuUsagePtr)) + 1
	if numGoroutines < 1 {
		numGoroutines = 1
	}

	fmt.Printf("Starting CPU stress with %d goroutines targeting %.2f CPU usage...\n", numGoroutines, *cpuUsagePtr)

	done := make(chan struct{})

	// Prevent channel close race condition
	var doneClosed sync.Once
	closeDone := func() {
		doneClosed.Do(func() {
			close(done)
		})
	}

	// Capture termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	var stopFlag int32

	// Launch worker goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Base workload parameters
			baseWorkload := 500 * time.Microsecond
			baseIdleTime := time.Duration((1 - *cpuUsagePtr) / *cpuUsagePtr * float64(baseWorkload))
			if baseIdleTime < 1*time.Microsecond {
				baseIdleTime = 1 * time.Microsecond
			}

			for {
				if atomic.LoadInt32(&stopFlag) == 1 {
					return
				}

				// Get current adjustment factor
				scaleFactor := monitor.AdjustWorkload()

				// Scale the workload
				workDuration := time.Duration(float64(baseWorkload) * scaleFactor)

				// Perform work (busy-wait)
				startWork := time.Now()
				for time.Since(startWork) < workDuration {
					// CPU-intensive operations
					_ = math.Sqrt(rand.Float64()) * math.Sqrt(rand.Float64())
				}

				// Calculate appropriate idle time based on desired duty cycle
				idleTime := time.Duration(float64(workDuration) * (1 - *cpuUsagePtr) / *cpuUsagePtr)
				if idleTime < 1*time.Microsecond {
					idleTime = 1 * time.Microsecond
				}

				// Idle for the calculated time
				time.Sleep(idleTime)
			}
		}()
	}

	// Start the monitoring goroutine
	go func() {
		// Wait for initial stabilization
		time.Sleep(500 * time.Millisecond)

		const monitorInterval = 1 * time.Second
		for {
			if atomic.LoadInt32(&stopFlag) == 1 {
				return
			}

			// Take a series of samples
			var totalUsage float64
			const numSamples = 3
			for i := 0; i < numSamples; i++ {
				start := time.Now()
				var counter uint64
				for j := 0; j < 100000 && time.Since(start) < 100*time.Millisecond; j++ {
					counter++
					_ = math.Sqrt(rand.Float64())
				}
				elapsed := time.Since(start)
				cpuEff := float64(counter) / float64(elapsed.Nanoseconds())
				// Calculate as a percentage of baseline
				currentUsage := cpuEff / baselineCPU
				totalUsage += currentUsage
				time.Sleep(10 * time.Millisecond)
			}

			// Calculate average and adjust for number of CPUs
			avgUsage := (totalUsage / numSamples) / float64(numCPU)

			// Update CPU usage monitor with observed usage
			monitor.UpdateUsage(avgUsage)

			time.Sleep(monitorInterval)
		}
	}()

	go func() {
		// Wait for termination signal
		<-quit
		fmt.Println("\nTermination signal received. Stopping CPU stress...")
		atomic.StoreInt32(&stopFlag, 1)
		closeDone()
	}()

	if !*runForeverPtr {
		time.Sleep(*durationPtr)
		fmt.Println("\nCPU stress completed.")
		atomic.StoreInt32(&stopFlag, 1)
		closeDone()
		// Keep the process running to prevent the pod from restarting
		select {}
	} else {
		// Run stress indefinitely
		fmt.Println("CPU stress will run indefinitely. Press Ctrl+C to stop.")
		<-done
	}
}
