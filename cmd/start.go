package cmd

// ... other imports

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	//	"os"
	//	"os/signal"
	"regexp"
	"strconv"
	"time"
)

// ... other code

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a session or break",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: Please specify 'session' or 'break'")
			return
		}

		duration, err := parseDuration(args[1:])
		if err != nil {
			fmt.Println(err)
			return
		}

		startTimer(args[0], duration)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func parseDuration(args []string) (time.Duration, error) {
	if len(args) == 0 {
		return 0, errors.New("missing duration argument")
	}

	input := args[0]

	// 1. Check for simple minutes:
	if minutes, err := strconv.Atoi(input); err == nil {
		return time.Duration(minutes) * time.Minute, nil
	}

	// 2. Check for "until" format:
	untilRegex := regexp.MustCompile(`^until\s+(?P<hour>\d{1,2})(?P<minute>\d{2})$`)
	match := untilRegex.FindStringSubmatch(input)

	if match != nil {
		now := time.Now()
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])

		targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if targetTime.Before(now) {
			targetTime = targetTime.Add(24 * time.Hour) // Target time is tomorrow
		}

		return targetTime.Sub(now), nil
	}

	// If none of the above match...
	return 0, errors.New("invalid duration format")
}

func startTimer(mode string, duration time.Duration) {
	bar := progressbar.NewOptions64( // Adjust total for precision if needed
		int64(duration.Seconds()),
		progressbar.OptionSetPredictTime(false),             // Disable automatic ETA
		progressbar.OptionShowBytes(false),                  // Display remaining time
		progressbar.OptionSetWidth(20),                      // Adjust progress bar width
		progressbar.OptionSetDescription("Time Remaining:"), // Optional
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	// Countdown loop
	for i := 0; i < int(duration.Seconds()); i++ {
		bar.Describe(fmt.Sprintf("%d seconds remaining", int(duration.Seconds())-i)) // Updated
		bar.Add(1)
		time.Sleep(time.Second)
	}

	playFinishedSound(mode)
}

func playFinishedSound(mode string) {
	// Implement playing a sound (you'll need a suitable library)
	fmt.Println("Sound not implemented")
}
