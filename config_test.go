package mongotest_test

import (
	"testing"

	"github.com/pinzolo/mongotest"
)

func TestDefaultConfig(t *testing.T) {
	defer mongotest.DefaultConfig()()
	if got := mongotest.GetConfigURL(); got != "" {
		t.Errorf("default url configuration should be empty but got %q", got)
	}
	if got := mongotest.GetConfigDatabase(); got != "" {
		t.Errorf("default database name configuration should be empty but got %q", got)
	}
	if got := mongotest.GetConfigFixtureRootDir(); got != "" {
		t.Errorf("default root directory of fixtures configuration should be empty but got %q", got)
	}
	if got := mongotest.GetConfigTimeout(); got != mongotest.DefaultTimeoutSeconds {
		t.Errorf("default timeoute seconds is invalid(want %d, got %d)", mongotest.DefaultTimeoutSeconds, got)
	}
	if got := mongotest.GetConfigFixtureFormat(); got != mongotest.FixtureFormatAuto {
		t.Errorf("default fixture format configuration should be auto but got %q", got)
	}
}

func TestURL(t *testing.T) {
	defer mongotest.DefaultConfig()()
	want := "mongodb://localhost:27017"
	mongotest.Configure(mongotest.URL(want))
	if got := mongotest.GetConfigURL(); got != want {
		t.Errorf("url configuration is invalid. (want: %q, got: %q)", want, got)
	}
}

func TestDatabase(t *testing.T) {
	defer mongotest.DefaultConfig()()
	want := "users"
	mongotest.Configure(mongotest.Database(want))
	if got := mongotest.GetConfigDatabase(); got != want {
		t.Errorf("database name configuration is invalid. (want: %q, got: %q)", want, got)
	}
}

func TestFixtureRootDir(t *testing.T) {
	defer mongotest.DefaultConfig()()
	want := "testdata"
	mongotest.Configure(mongotest.FixtureRootDir(want))
	if got := mongotest.GetConfigFixtureRootDir(); got != want {
		t.Errorf("fixture root directory configuration is invalid. (want: %q, got: %q)", want, got)
	}
}

func TestFixtureFormat(t *testing.T) {
	defer mongotest.DefaultConfig()()
	want := mongotest.FixtureFormatJSON
	mongotest.Configure(mongotest.FixtureFormat(want))
	if got := mongotest.GetConfigFixtureFormat(); got != want {
		t.Errorf("fixture format configuration is invalid. (want: %q, got: %q)", want, got)
	}
}

func TestTimeout(t *testing.T) {
	defer mongotest.DefaultConfig()()
	want := 30
	mongotest.Configure(mongotest.Timeout(30))
	if got := mongotest.GetConfigTimeout(); got != want {
		t.Errorf("timeout seconds configuration is invalid. (want: %d, got: %d)", want, got)
	}
}
