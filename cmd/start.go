package cmd

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
	Run:   startFocusSessionCommand,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

type SessionMode string

const (
	FocusSession SessionMode = "focus"
	BreakSession SessionMode = "break"
)

func startFocusSessionCommand(command *cobra.Command, commandLineArguments []string) {
	duration, parseError := parseDurationFromArguments(commandLineArguments)
	if parseError != nil {
		log.Fatal(parseError)
	}
	startSession(FocusSession, duration)
}

func parseDurationFromArguments(arguments []string) (time.Duration, error) {
	if len(arguments) == 0 {
		return 0, nil
	}
	durationInput := arguments[0]
	minutes, conversionError := strconv.Atoi(durationInput)
	if conversionError != nil {
		return 0, errors.New("invalid duration format")
	}
	return time.Duration(minutes) * time.Minute, nil
}

// Passing a session duration of 0 will result in using the default session duration as defined in the config
func startSession(mode SessionMode, duration time.Duration) {
	if duration == 0 {
		duration = getDefaultDuration(mode)
	}
	fmt.Printf("Starting timer for %v\n", duration)
	startTimer(mode, duration)
}

func getDefaultDuration(mode SessionMode) time.Duration {
	var defaultDuration int
	if mode == FocusSession {
		defaultDuration = viper.GetInt(configuration.DefaultSessionDuration)
	} else if mode == BreakSession {
		defaultDuration = viper.GetInt(configuration.DefaultBreakDuration)
	}
	return time.Duration(defaultDuration) * time.Minute
}

func startTimer(mode SessionMode, duration time.Duration) {
	bar := createStatusBar(duration)
	beginCountdown(duration, bar)
	playFinishedSound()
	endSession(mode)
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

func endSession(mode SessionMode) {
	var choice string
	if mode == FocusSession {
		fmt.Print("Start a [b]reak, or e[x]it?\n")
		fmt.Scanln(&choice)

		switch choice {
		case "b":
			startSession(BreakSession, 0)
		case "x":
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, Exiting the program...")
			os.Exit(0)
		}
	} else if mode == BreakSession {
		fmt.Print("Start a [s]ession, or e[x]it?\n")
		fmt.Scanln(&choice)

		switch choice {
		case "s":
			startSession(FocusSession, 0)
		case "x":
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, Exiting the program...")
			os.Exit(0)
		}
	}

}
