package main

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCgroupV2CPUMax(t *testing.T) {
	tests := []struct {
		content string
		want    float64
		wantOK  bool
	}{
		{"200000 100000\n", 2, true},  // limits.cpu: 2
		{"50000 100000\n", 0.5, true}, // limits.cpu: 500m
		{"20000 100000", 0.2, true},   // limits.cpu: 200m
		{"max 100000\n", 0, false},    // no limit
		{"", 0, false},
		{"garbage", 0, false},
	}
	for _, tt := range tests {
		got, ok := parseCgroupV2CPUMax(tt.content)
		if ok != tt.wantOK || math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("parseCgroupV2CPUMax(%q) = (%v, %v), want (%v, %v)", tt.content, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestParseCgroupV1Quota(t *testing.T) {
	tests := []struct {
		quota, period string
		want          float64
		wantOK        bool
	}{
		{"150000\n", "100000\n", 1.5, true},
		{"-1\n", "100000\n", 0, false}, // no limit
		{"x", "100000", 0, false},
	}
	for _, tt := range tests {
		got, ok := parseCgroupV1Quota(tt.quota, tt.period)
		if ok != tt.wantOK || math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("parseCgroupV1Quota(%q, %q) = (%v, %v), want (%v, %v)", tt.quota, tt.period, got, ok, tt.want, tt.wantOK)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectCPULimitIn(t *testing.T) {
	v2 := t.TempDir()
	writeFile(t, filepath.Join(v2, "cpu.max"), "50000 100000\n")
	if got, ok := detectCPULimitIn(v2); !ok || got != 0.5 {
		t.Errorf("cgroup v2: got (%v, %v), want (0.5, true)", got, ok)
	}

	v2NoLimit := t.TempDir()
	writeFile(t, filepath.Join(v2NoLimit, "cpu.max"), "max 100000\n")
	if _, ok := detectCPULimitIn(v2NoLimit); ok {
		t.Error("cgroup v2 without limit: want ok=false")
	}

	v1 := t.TempDir()
	writeFile(t, filepath.Join(v1, "cpu,cpuacct", "cpu.cfs_quota_us"), "200000\n")
	writeFile(t, filepath.Join(v1, "cpu,cpuacct", "cpu.cfs_period_us"), "100000\n")
	if got, ok := detectCPULimitIn(v1); !ok || got != 2 {
		t.Errorf("cgroup v1: got (%v, %v), want (2, true)", got, ok)
	}

	if _, ok := detectCPULimitIn(t.TempDir()); ok {
		t.Error("no cgroup files: want ok=false")
	}
}
