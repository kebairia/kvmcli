package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is a global logger instance that other packages can use.
var Log = logrus.New()

func init() {
	// Configure Logrus to output colored logs with full timestamps.
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true, // This pads the level text for better spacing.
		// Timestamp configuration
		DisableTimestamp: true, // Set to true to disable timestamps
		FullTimestamp:    true,
		TimestampFormat:  "2006/01/02 15:04:05",
	})

	// Set the default log level (you can adjust this as needed).
	Log.SetLevel(logrus.DebugLevel)
}
