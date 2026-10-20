package main

import (
	"math"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestCoresPlan(t *testing.T) {
	tests := []struct {
		cores       float64
		wantWorkers int
		wantDuty    float64
	}{
		{0.1, 1, 0.1},
		{0.2, 1, 0.2},
		{1, 1, 1},
		{1.5, 2, 0.75},
		{2, 2, 1},
		{3.3, 4, 0.825},
	}
	for _, tt := range tests {
		workers, duty := coresPlan(tt.cores)
		if workers != tt.wantWorkers || math.Abs(duty-tt.wantDuty) > 1e-9 {
			t.Errorf("coresPlan(%v) = (%d, %v), want (%d, %v)", tt.cores, workers, duty, tt.wantWorkers, tt.wantDuty)
		}
		if got := float64(workers) * duty; math.Abs(got-tt.cores) > 1e-9 {
			t.Errorf("coresPlan(%v): workers*duty = %v, want %v", tt.cores, got, tt.cores)
		}
		if duty <= 0 || duty > 1 {
			t.Errorf("coresPlan(%v): duty %v out of (0, 1]", tt.cores, duty)
		}
	}
}

// The -cpu worker count must not change: existing users depend on the load it produces.
func TestLegacyWorkerCountUnchanged(t *testing.T) {
	tests := []struct {
		numCPU int
		cpu    float64
		want   int
	}{
		{1, 0.2, 1},
		{1, 1.0, 2},
		{4, 0.2, 1},
		{14, 0.2, 3},
		{14, 0.5, 8},
		{64, 0.5, 33},
	}
	for _, tt := range tests {
		if got := legacyWorkerCount(tt.numCPU, tt.cpu); got != tt.want {
			t.Errorf("legacyWorkerCount(%d, %v) = %d, want %d", tt.numCPU, tt.cpu, got, tt.want)
		}
	}
}

func TestNextCorrection(t *testing.T) {
	const maxC = 4.0
	tests := []struct {
		name                       string
		correction, target, actual float64
		want                       float64
	}{
		{"on target", 1, 0.5, 0.5, 1},
		{"no measurement", 1.2, 0.5, 0, 1.2},
		{"slightly low", 1, 0.5, 0.45, 0.5 / 0.45},
		{"far too low is step-limited", 1, 0.5, 0.1, maxCorrectionStep},
		{"far too high is step-limited", 1, 0.5, 2, 1 / maxCorrectionStep},
		{"clamped to max", 3.9, 0.5, 0.1, maxC},
		{"clamped to min", 0.11, 0.5, 5, minCorrection},
	}
	for _, tt := range tests {
		got := nextCorrection(tt.correction, tt.target, tt.actual, maxC)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("%s: nextCorrection(%v, %v, %v) = %v, want %v", tt.name, tt.correction, tt.target, tt.actual, got, tt.want)
		}
	}
}

// Runs real load and checks the measured process CPU against the -cores target.
func TestCoresLoadAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("load test skipped in -short mode")
	}
	for _, cores := range []float64{0.2, 0.5, 1.5} {
		if int(math.Ceil(cores)) > runtime.NumCPU()-1 {
			t.Logf("skipping %.1f cores: only %d CPUs", cores, runtime.NumCPU())
			continue
		}
		c := cores
		var stop int32
		startCoresStress(func() float64 { return c }, &stop)
		time.Sleep(3 * time.Second) // let the controller settle

		actual := measureCores(3 * time.Second)
		atomic.StoreInt32(&stop, 1)
		time.Sleep(1200 * time.Millisecond) // let workers and controller exit

		checkCores(t, "-cores", cores, actual)
	}
}

// The target can change while running (in-place pod resize with -limit-fraction),
// including growing beyond the initial worker count.
func TestCoresTargetChange(t *testing.T) {
	if testing.Short() {
		t.Skip("load test skipped in -short mode")
	}
	if runtime.NumCPU() < 4 {
		t.Skip("needs at least 4 CPUs")
	}
	var target uint64
	var stop int32
	atomic.StoreUint64(&target, math.Float64bits(0.4))
	startCoresStress(func() float64 { return math.Float64frombits(atomic.LoadUint64(&target)) }, &stop)

	time.Sleep(3 * time.Second)
	checkCores(t, "start", 0.4, measureCores(2*time.Second))

	atomic.StoreUint64(&target, math.Float64bits(1.6)) // grow: needs a 2nd worker
	time.Sleep(4 * time.Second)
	checkCores(t, "after grow", 1.6, measureCores(2*time.Second))

	atomic.StoreUint64(&target, math.Float64bits(0.25)) // shrink
	time.Sleep(4 * time.Second)
	checkCores(t, "after shrink", 0.25, measureCores(2*time.Second))

	atomic.StoreInt32(&stop, 1)
	time.Sleep(1200 * time.Millisecond)
}

func measureCores(d time.Duration) float64 {
	cpu0, wall0 := processCPUTime(), time.Now()
	time.Sleep(d)
	return float64(processCPUTime()-cpu0) / float64(time.Since(wall0))
}

func checkCores(t *testing.T, label string, want, got float64) {
	t.Helper()
	if tol := 0.1*want + 0.05; math.Abs(got-want) > tol {
		t.Errorf("%s: measured %.3f cores, want %.2f ±%.2f", label, got, want, tol)
	} else {
		t.Logf("%s: measured %.3f cores (target %.2f)", label, got, want)
	}
}
