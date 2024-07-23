package cmd

// ... other imports

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Justin-Arnold/epoch-cli/internal/configuration"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed sounds/yeahboi.mp3
var yeahboi []byte

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a session or break",
	Run:   startSession,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startSession(command *cobra.Command, commandLineArguments []string) {
	duration, getDurationError := getDuration(commandLineArguments)
	if getDurationError != nil {
		log.Fatal(getDurationError)
	}
	fmt.Printf("Starting timer for %v\n", duration)
	startTimer(duration)
}

func getDuration(args []string) (time.Duration, error) {
	defaultDuration := viper.GetInt(configuration.ConfigOptionSessionDuration)
	if len(args) == 0 {
		return time.Duration(defaultDuration) * time.Minute, nil
	} else {
		return parseDuration(args)
	}
}

func parseDuration(args []string) (time.Duration, error) {
	durationInput := args[0]
	if minutes, err := strconv.Atoi(durationInput); err == nil {
		return time.Duration(minutes) * time.Minute, nil
	}
	return 0, errors.New("invalid duration format")
}

func startTimer(duration time.Duration) {
	bar := createStatusBar(duration)
	beginCountdown(duration, bar)
	playFinishedSound()
	endSession(duration)
}

func beginCountdown(duration time.Duration, statusBar *progressbar.ProgressBar) {
	totalSeconds := int(duration.Seconds())
	for i := 0; i < totalSeconds; i++ {
		remaining := duration - time.Duration(i)*time.Second
		statusBar.Describe(fmt.Sprintf("%s remaining", formatDuration(remaining)))
		statusBar.Add(1)
		time.Sleep(time.Second)
	}
}

func createStatusBar(duration time.Duration) *progressbar.ProgressBar {
	return progressbar.NewOptions64( // Adjust total for precision if needed
		int64(duration.Seconds()),
		progressbar.OptionSetPredictTime(false),             // Disable automatic ETA
		progressbar.OptionShowBytes(false),                  // Display remaining time
		progressbar.OptionSetWidth(20),                      // Adjust progress bar width
		progressbar.OptionSetDescription("Time Remaining:"), // Optional
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)
}

// Helper function to format duration
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func playFinishedSound() {
	streamer, format := getSound()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))
	<-done
}

func getSound() (beep.StreamSeekCloser, beep.Format) {
	reader := bytes.NewReader(yeahboi)
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	return streamer, format
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
		startTimer(time.Duration(originalDuration))
	case "n":
		fmt.Print("How many minutes should the new session be?\n")
		fmt.Scanln(&choice)
		inputDuration, err := strconv.Atoi(choice)
		if err != nil {
			log.Fatal("Invalid Duration, exiting the program")
		}
		startTimer(time.Duration(inputDuration) * time.Minute)
	case "x":
		fmt.Println("Exiting the program...")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice, Exiting the program...")
		os.Exit(0)
	}
}
