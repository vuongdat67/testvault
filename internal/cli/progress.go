package cli

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar represents a terminal progress bar
type ProgressBar struct {
	current   int64
	total     int64
	width     int
	operation string
	startTime time.Time
	lastPrint time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64, operation string) *ProgressBar {
	return &ProgressBar{
		current:   0,
		total:     total,
		width:     50,
		operation: operation,
		startTime: time.Now(),
		lastPrint: time.Now(),
	}
}

// Update updates the progress bar with current progress
func (pb *ProgressBar) Update(current int64) {
	pb.current = current
	
	// Only update display every 100ms to avoid too frequent updates
	now := time.Now()
	if now.Sub(pb.lastPrint) < 100*time.Millisecond && current < pb.total {
		return
	}
	pb.lastPrint = now
	
	pb.display()
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.current = pb.total
	pb.display()
	fmt.Println()
}

// display renders the progress bar
func (pb *ProgressBar) display() {
	if pb.total == 0 {
		return
	}
	
	// Calculate percentage
	percentage := float64(pb.current) / float64(pb.total) * 100
	
	// Calculate ETA
	elapsed := time.Since(pb.startTime)
	var eta string
	if pb.current > 0 {
		totalTime := time.Duration(float64(elapsed) * float64(pb.total) / float64(pb.current))
		remainingTime := totalTime - elapsed
		if remainingTime > 0 {
			eta = fmt.Sprintf(" ETA: %s", FormatDuration(remainingTime.Seconds()))
		}
	}
	
	// Calculate speed
	speed := ""
	if elapsed.Seconds() > 0 {
		bytesPerSecond := float64(pb.current) / elapsed.Seconds()
		speed = fmt.Sprintf(" %s/s", FormatBytes(uint64(bytesPerSecond)))
	}
	
	// Create progress bar
	filledWidth := int(float64(pb.width) * percentage / 100)
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", pb.width-filledWidth)
	
	// Format output
	output := fmt.Sprintf("\r%s [%s] %.1f%% %s/%s%s%s",
		pb.operation,
		bar,
		percentage,
		FormatBytes(uint64(pb.current)),
		FormatBytes(uint64(pb.total)),
		speed,
		eta,
	)
	
	fmt.Print(output)
}

// SimpleProgress shows a simple text-based progress update
type SimpleProgress struct {
	operation string
	startTime time.Time
	lastUpdate time.Time
}

// NewSimpleProgress creates a new simple progress tracker
func NewSimpleProgress(operation string) *SimpleProgress {
	return &SimpleProgress{
		operation:  operation,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Update shows progress update
func (sp *SimpleProgress) Update(current, total int64) {
	// Only update every second to avoid spam
	now := time.Now()
	if now.Sub(sp.lastUpdate) < time.Second && current < total {
		return
	}
	sp.lastUpdate = now
	
	if total > 0 {
		percentage := float64(current) / float64(total) * 100
		fmt.Printf("\r%s... %.1f%% (%s/%s)", 
			sp.operation,
			percentage,
			FormatBytes(uint64(current)),
			FormatBytes(uint64(total)),
		)
	} else {
		fmt.Printf("\r%s... %s", sp.operation, FormatBytes(uint64(current)))
	}
}

// Finish completes the progress tracking
func (sp *SimpleProgress) Finish() {
	elapsed := time.Since(sp.startTime)
	fmt.Printf(" ✅ Done (%s)\n", FormatDuration(elapsed.Seconds()))
}

// SpinnerProgress shows a spinning progress indicator
type SpinnerProgress struct {
	operation string
	frames    []string
	current   int
	ticker    *time.Ticker
	done      chan bool
}

// NewSpinnerProgress creates a new spinner progress indicator
func NewSpinnerProgress(operation string) *SpinnerProgress {
	return &SpinnerProgress{
		operation: operation,
		frames:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		current:   0,
		done:      make(chan bool),
	}
}

// Start starts the spinner
func (sp *SpinnerProgress) Start() {
	sp.ticker = time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-sp.done:
				return
			case <-sp.ticker.C:
				fmt.Printf("\r%s %s...", sp.frames[sp.current%len(sp.frames)], sp.operation)
				sp.current++
			}
		}
	}()
}

// Stop stops the spinner
func (sp *SpinnerProgress) Stop() {
	if sp.ticker != nil {
		sp.ticker.Stop()
	}
	sp.done <- true
	fmt.Print("\r")
}

// ProgressCallback is a function type for progress callbacks
type ProgressCallback func(current, total int64, operation string)

// NoOpProgress is a no-operation progress callback
func NoOpProgress(current, total int64, operation string) {
	// Do nothing - for quiet mode
}

// ConsoleProgress is a console-based progress callback
func ConsoleProgress(current, total int64, operation string) {
	if total > 0 {
		percentage := float64(current) / float64(total) * 100
		fmt.Printf("\rProgress: %.1f%% (%s/%s)",
			percentage,
			FormatBytes(uint64(current)),
			FormatBytes(uint64(total)),
		)
		if current >= total {
			fmt.Println(" ✅")
		}
	}
}

// DetailedProgress provides detailed progress information
func DetailedProgress(current, total int64, operation string) {
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("[%.1f%%] %s: %s/%s\n",
		percentage,
		operation,
		FormatBytes(uint64(current)),
		FormatBytes(uint64(total)),
	)
}
