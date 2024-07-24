package configuration

const (
	DefaultFocusDuration ConfigOptionKey = "default_session_duration"
	DefaultBreakDuration ConfigOptionKey = "default_break_duration"
)

var DefaultDurationOptions = ConfigPartial{
	DefaultFocusDuration: 25,
	DefaultBreakDuration: 5,
}

func init() {
	RegisterConfigOptions(DefaultDurationOptions)
}
