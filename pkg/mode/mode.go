package mode

import "os"

const envAppMode = "APP_MODE"

const (
	// DebugMode indicates app mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates app mode is release.
	ReleaseMode = "release"
	// TestMode indicates app mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

var (
	appMode  = debugCode
	modeName = DebugMode
)

func init() {
	mode := os.Getenv(envAppMode)
	SetMode(mode)
}

// SetMode sets mode according to input string.
func SetMode(value string) {
	if value == "" {
		value = DebugMode
	}

	switch value {
	case DebugMode:
		appMode = debugCode
	case ReleaseMode:
		appMode = releaseCode
	case TestMode:
		appMode = testCode
	default:
		panic("mode unknown: " + value)
	}

	modeName = value
}

// Mode returns current app mode.
func Mode() string {
	return modeName
}
