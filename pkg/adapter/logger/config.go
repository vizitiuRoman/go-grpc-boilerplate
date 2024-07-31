package logger

type Config struct {
	// Level sets the logger's level. Possible values are "debug", "info", "warn", "error".
	Level string `json:"level" yaml:"level"`
	// Encoding sets the logger's encoding.Possible values are "json" and "console".
	Encoding string `json:"encoding" yaml:"encoding"`
	// LevelEncoder is the level encoder type. Possible values are "capital", "capitalColor", "color", "lower".
	LevelEncoder string `json:"level_encoder" yaml:"level_encoder"`
	// CallerSkip is the number of stack frames to skip to find the caller.
	CallerSkip *int `json:"caller_skip" yaml:"caller_skip"`
}
