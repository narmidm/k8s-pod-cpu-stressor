package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"syscall"
	"time"
)

// Precise mode (-cores): generates exactly N cores of load, independent of the
// node size, and corrects the duty cycle using this process's measured CPU time.

const (
	coresPeriod          = 10 * time.Millisecond // one busy+idle cycle per worker
	coresControlInterval = 1 * time.Second
	minCorrection        = 0.1
	maxCorrectionStep    = 1.25 // limit change per control interval to avoid oscillation
)

// coresPlan returns the number of workers and the busy fraction of each,
// so that workers*duty == cores.
func coresPlan(cores float64) (workers int, duty float64) {
	workers = int(math.Ceil(cores))
	if workers < 1 {
		workers = 1
	}
	return workers, cores / float64(workers)
}

// nextCorrection adjusts the duty-cycle multiplier towards target/actual,
// limited per step and clamped to [minCorrection, maxCorrection].
func nextCorrection(correction, target, actual, maxCorrection float64) float64 {
	if actual <= 0 || target <= 0 {
		return correction
	}
	ratio := target / actual
	if ratio > maxCorrectionStep {
		ratio = maxCorrectionStep
	} else if ratio < 1/maxCorrectionStep {
		ratio = 1 / maxCorrectionStep
	}
	c := correction * ratio
	if c > maxCorrection {
		c = maxCorrection
	}
	if c < minCorrection {
		c = minCorrection
	}
	return c
}

// processCPUTime returns the user+system CPU time consumed by this process.
func processCPUTime() time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

// runCoresWorker is worker number index. The current load (in cores) is spread
// evenly over the first ceil(load) workers; any others stay idle.
func runCoresWorker(index int, loadBits *uint64, stopFlag *int32) {
	for atomic.LoadInt32(stopFlag) == 0 {
		active, duty := coresPlan(math.Float64frombits(atomic.LoadUint64(loadBits)))
		if index >= active {
			time.Sleep(coresPeriod)
			continue
		}
		busy := time.Duration(float64(coresPeriod) * duty)

		start := time.Now()
		for time.Since(start) < busy {
			_ = math.Sqrt(rand.Float64()) * math.Sqrt(rand.Float64())
		}

		// Sleep until the end of the period, so oversleeping doesn't accumulate
		if idle := coresPeriod - time.Since(start); idle > 0 {
			time.Sleep(idle)
		}
	}
}

// warnIfOverLimit prints a warning when the target can't be reached within the CPU limit.
func warnIfOverLimit(cores float64) {
	if avail, source := availableCPUs(); cores > avail {
		fmt.Printf("Warning: target %.2f cores exceeds the %s of %.2f cores; load will be capped\n", cores, source, avail)
	}
}

// startCoresStress launches the workers and the feedback controller.
// target is re-evaluated every control interval, so a target derived from the
// CPU limit follows in-place pod resizes. It returns immediately; everything
// stops once stopFlag is set to 1.
func startCoresStress(target func() float64, stopFlag *int32) {
	cores := target()
	workers, _ := coresPlan(cores)
	warnIfOverLimit(cores)
	fmt.Printf("Starting CPU stress with %d workers targeting %.2f cores...\n", workers, cores)

	var loadBits uint64
	atomic.StoreUint64(&loadBits, math.Float64bits(cores))
	spawned := 0
	spawn := func(n int) {
		for ; spawned < n; spawned++ {
			go runCoresWorker(spawned, &loadBits, stopFlag)
		}
	}
	spawn(workers)

	go func() {
		correction := 1.0
		lastCPU, lastWall := processCPUTime(), time.Now()
		for {
			time.Sleep(coresControlInterval)
			if atomic.LoadInt32(stopFlag) == 1 {
				return
			}

			if t := target(); math.Abs(t-cores) > 1e-6 {
				fmt.Printf("Target changed: %.2f -> %.2f cores\n", cores, t)
				cores = t
				warnIfOverLimit(cores)
			}

			cpu, wall := processCPUTime(), time.Now()
			actual := float64(cpu-lastCPU) / float64(wall.Sub(lastWall))
			lastCPU, lastWall = cpu, wall

			// Never generate more than ceil(target) cores (all workers fully busy)
			maxCorrection := math.Ceil(cores) / cores
			correction = nextCorrection(correction, cores, actual, maxCorrection)
			load := cores * correction
			n, duty := coresPlan(load)
			spawn(n)
			atomic.StoreUint64(&loadBits, math.Float64bits(load))

			fmt.Printf("CPU Usage: %.2f cores (target: %.2f cores), workers: %d, duty: %.0f%%\n",
				actual, cores, n, duty*100)
		}
	}()
}
