package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Container CPU limit detection from cgroup files (Kubernetes resources.limits.cpu).

const cgroupRoot = "/sys/fs/cgroup"

// parseCgroupV2CPUMax parses cgroup v2 cpu.max ("<quota> <period>" or "max <period>").
// ok is false when there is no limit or the content is invalid.
func parseCgroupV2CPUMax(content string) (cores float64, ok bool) {
	fields := strings.Fields(content)
	if len(fields) != 2 || fields[0] == "max" {
		return 0, false
	}
	quota, err1 := strconv.ParseFloat(fields[0], 64)
	period, err2 := strconv.ParseFloat(fields[1], 64)
	if err1 != nil || err2 != nil || quota <= 0 || period <= 0 {
		return 0, false
	}
	return quota / period, true
}

// parseCgroupV1Quota parses cgroup v1 cpu.cfs_quota_us and cpu.cfs_period_us (-1 quota = no limit).
func parseCgroupV1Quota(quotaContent, periodContent string) (cores float64, ok bool) {
	quota, err1 := strconv.ParseFloat(strings.TrimSpace(quotaContent), 64)
	period, err2 := strconv.ParseFloat(strings.TrimSpace(periodContent), 64)
	if err1 != nil || err2 != nil || quota <= 0 || period <= 0 {
		return 0, false
	}
	return quota / period, true
}

// detectCPULimitIn looks for a CPU limit under root (cgroup v2 first, then v1).
func detectCPULimitIn(root string) (cores float64, ok bool) {
	if b, err := os.ReadFile(filepath.Join(root, "cpu.max")); err == nil {
		return parseCgroupV2CPUMax(string(b))
	}
	for _, dir := range []string{"cpu", "cpu,cpuacct", "cpuacct,cpu"} {
		q, err1 := os.ReadFile(filepath.Join(root, dir, "cpu.cfs_quota_us"))
		p, err2 := os.ReadFile(filepath.Join(root, dir, "cpu.cfs_period_us"))
		if err1 == nil && err2 == nil {
			return parseCgroupV1Quota(string(q), string(p))
		}
	}
	return 0, false
}

// availableCPUs returns the container's CPU limit in cores if one is set,
// otherwise the number of CPUs visible to the process.
func availableCPUs() (cores float64, source string) {
	if c, ok := detectCPULimitIn(cgroupRoot); ok {
		return c, "container CPU limit"
	}
	return float64(runtime.NumCPU()), "node CPUs (no container CPU limit)"
}
