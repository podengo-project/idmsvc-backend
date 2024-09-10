package perf

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessData represents the data collected for a process.
type ProcessData struct {
	CPUPercent float64
	MemRSS     uint64
	Time       time.Time
}

// ProcessMonitor monitors a specific process and stores its data.
type ProcessMonitor struct {
	process            *process.Process
	processName        string
	running            bool
	readInterval       time.Duration
	maxHistoryDuration time.Duration
	mutex              sync.Mutex
	history            []ProcessData
}

// NewProcessMonitor creates a new ProcessMonitor instance.
func NewProcessMonitor(
	processName string, readInterval, maxHistoryDuration time.Duration,
) (*ProcessMonitor, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	process := getProcess(processes, processName)
	if process == nil {
		return nil, fmt.Errorf("process %s not found", processName)
	}

	return &ProcessMonitor{
		processName:        processName,
		process:            process,
		readInterval:       readInterval,
		maxHistoryDuration: maxHistoryDuration,
		history:            make([]ProcessData, 0),
		running:            false,
		mutex:              sync.Mutex{},
	}, nil
}

// Start starts the monitoring process.
func (pm *ProcessMonitor) Start() {
	pm.running = true
	go func() {
		for {
			if !pm.running {
				break
			}
			pm.Monitor()
			time.Sleep(pm.readInterval)
		}
	}()
}

// Stop stops the monitoring process.
func (pm *ProcessMonitor) Stop() {
	pm.running = false
}

// GetHistory returns the history of collected data.
func (pm *ProcessMonitor) GetHistory() []ProcessData {
	return pm.history
}

// ClearHistory clears the history of collected data.
func (pm *ProcessMonitor) ClearHistory() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.history = make([]ProcessData, 0)
}

// MaxCPUPercent returns the maximum CPU usage.
func (pm *ProcessMonitor) MaxCPUPercent() float64 {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var max float64
	for _, data := range pm.history {
		if data.CPUPercent > max {
			max = data.CPUPercent
		}
	}
	return max
}

// MaxMemRSS returns the maximum memory usage.
func (pm *ProcessMonitor) MaxMemRSS() uint64 {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var max uint64
	for _, data := range pm.history {
		if data.MemRSS > max {
			max = data.MemRSS
		}
	}
	return max
}

// MaxCPUPercentInTimeRange returns the peak CPU usage in a specific time range.
func (pm *ProcessMonitor) MaxCPUPercentInTimeRange(start, end time.Time) float64 {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var peak float64
	for _, data := range pm.history {
		if data.Time.After(start) && data.Time.Before(end) {
			if data.CPUPercent > peak {
				peak = data.CPUPercent
			}
		}
	}
	return peak
}

// MaxMemRSSInTimeRange returns the peak memory usage in a specific time range.
func (pm *ProcessMonitor) MaxMemRSSInTimeRange(start, end time.Time) uint64 {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var peak uint64
	for _, data := range pm.history {
		if data.Time.After(start) && data.Time.Before(end) {
			if data.MemRSS > peak {
				peak = data.MemRSS
			}
		}
	}
	return peak
}

// Monitor takes a sample of the process data.
func (pm *ProcessMonitor) Monitor() {
	cpuPercent, cpuErr := pm.process.CPUPercent()

	memInfo, memErr := pm.process.MemoryInfo()

	if cpuErr != nil || memErr != nil {
		return
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.history = append(pm.history, ProcessData{
		CPUPercent: cpuPercent,
		MemRSS:     memInfo.RSS,
		Time:       time.Now(),
	})

	pm.pruneHistory()
}

func getProcess(processes []*process.Process, name string) *process.Process {
	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}
		if strings.HasPrefix(cmdline, name) {
			return p
		}
	}
	return nil
}

func (pm *ProcessMonitor) pruneHistory() {
	if len(pm.history) == 0 {
		return
	}
	now := time.Now()
	for i, data := range pm.history {
		if now.Sub(data.Time) <= pm.maxHistoryDuration {
			pm.history = pm.history[i:]
			return
		}
	}
}

// HumanBytes converts bytes to human-readable format.
func HumanBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
