package cmd

// ... other imports

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

//go:embed sounds/yeahboi.mp3
var yeahboi []byte

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a session or break",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Error: missing duration argument")
			return
		}

		duration, err := parseDuration(args[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		startTimer("session", duration)
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
	endSession(duration)
}

func playFinishedSound(mode string) {
	reader := bytes.NewReader(yeahboi)
	streamer, format, err := mp3.Decode(ioutil.NopCloser(reader))
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))
	<-done
}

type OriginalSessionType string

const (
	DefaultTypeSession OriginalSessionType = "default"
	CustomTypeSession  OriginalSessionType = "custom"
)

func endSession(originalDuration time.Duration) {

	fmt.Print("Start another [s]ession, start a session with a [n]ew time, or e[x]it?\n")
	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "s":
		startTimer("session", time.Duration(originalDuration))
	case "n":
		fmt.Print("How many minutes should the new session be?\n")
		fmt.Scanln(&choice)
		inputDuration, err := strconv.Atoi(choice)
		if err != nil {
			log.Fatal("Invalid Duration, exiting the program")
		}
		startTimer("session", time.Duration(inputDuration)*time.Minute)
	case "x":
		fmt.Println("Exiting the program...")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice, Exiting the program...")
		os.Exit(0)
	}
}
