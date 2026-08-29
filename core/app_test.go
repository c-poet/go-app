package core

import (
	"log/slog"
	"testing"
)

func TestInitAppLogName(t *testing.T) {
	tests := []struct {
		name       string
		opts       []AppOption
		confName   string
		logName    string
		logEnabled bool
	}{
		{
			name:       "defaults to application name",
			opts:       []AppOption{WithName("service")},
			confName:   "service",
			logName:    "service",
			logEnabled: true,
		},
		{
			name:       "uses configured names",
			opts:       []AppOption{WithName("service"), WithConfName("service-config"), WithLogName("service-api"), WithLogEnabled(false)},
			confName:   "service-config",
			logName:    "service-api",
			logEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAppOptions(tt.opts...).LogName; got != tt.logName {
				t.Errorf("log name = %q, want %q", got, tt.logName)
			}
			if got := newAppOptions(tt.opts...).ConfName; got != tt.confName {
				t.Errorf("configuration name = %q, want %q", got, tt.confName)
			}
			if got := newAppOptions(tt.opts...).LogEnabled; got != tt.logEnabled {
				t.Errorf("log enabled = %t, want %t", got, tt.logEnabled)
			}
		})
	}
}

func TestReloadLogConfWhenDisabled(t *testing.T) {
	originalLogger := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(originalLogger)
	})

	app := &Application{logEnabled: false}
	app.ReloadLogConf()

	if slog.Default() != originalLogger {
		t.Error("ReloadLogConf() changed the default logger while logging is disabled")
	}
}
