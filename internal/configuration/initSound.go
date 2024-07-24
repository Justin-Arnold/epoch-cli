package configuration

const (
	CustomFocusSound ConfigOptionKey = "custom_focus_sound"
	CustomBreakSound ConfigOptionKey = "custom_break_sound"
	FocusSoundType   ConfigOptionKey = "focus_sound_type"
	BreakSoundType   ConfigOptionKey = "break_sound_type"
)

var DefaultSoundOptions = ConfigPartial{
	CustomFocusSound: "",
	CustomBreakSound: "",
	FocusSoundType:   "default",
	BreakSoundType:   "default",
}

func init() {
	RegisterConfigOptions(DefaultSoundOptions)

}
