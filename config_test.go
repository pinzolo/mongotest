package mongotest_test

import (
	"testing"

	"github.com/pinzolo/mongotest"
)

func TestDefaultConfig(t *testing.T) {
	defer mongotest.DefaultConfig()()
	conf := mongotest.Configuration()
	if got := conf.URL; got != "" {
		t.Errorf("default URL configuration should be empty but got %q", got)
	}
	if got := conf.Database; got != "" {
		t.Errorf("default Database name configuration should be empty but got %q", got)
	}
	if got := conf.FixtureRootDir; got != "" {
		t.Errorf("default root directory of fixtures configuration should be empty but got %q", got)
	}
	if got := conf.Timeout; got != mongotest.DefaultTimeoutSeconds {
		t.Errorf("default timeoute seconds is invalid(want %d, got %d)", mongotest.DefaultTimeoutSeconds, got)
	}
	if got := conf.FixtureFormat; got != mongotest.FixtureFormatAuto {
		t.Errorf("default fixture format configuration should be auto but got %q", got)
	}
}

func TestConfigure(t *testing.T) {
	defer mongotest.DefaultConfig()()
	c := mongotest.Config{
		URL:            "mongodb://localhost:27017",
		Database:       "mongotest",
		FixtureRootDir: "testdata",
		FixtureFormat:  mongotest.FixtureFormatJSON,
		Timeout:        30,
		PreInsertFuncs: []mongotest.PreInsertFunc{
			mongotest.SimpleConvertTime("users", "created_at"),
		},
	}
	mongotest.Configure(c)
	conf := mongotest.Configuration()
	if conf.URL != c.URL {
		t.Errorf("URL should be overwritten. (want: %q, got: %q", c.URL, conf.URL)
	}
	if conf.Database != c.Database {
		t.Errorf("Database should be overwritten. (want: %q, got: %q", c.Database, conf.Database)
	}
	if conf.FixtureRootDir != c.FixtureRootDir {
		t.Errorf("FixtureRootDir should be overwritten. (want: %q, got: %q", c.FixtureRootDir, conf.FixtureRootDir)
	}
	if conf.FixtureFormat != c.FixtureFormat {
		t.Errorf("FixtureFormat should be overwritten. (want: %q, got: %q", c.FixtureFormat, conf.FixtureFormat)
	}
	if conf.Timeout != c.Timeout {
		t.Errorf("Timeout should be overwritten. (want: %d, got: %d", c.Timeout, conf.Timeout)
	}
	if len(conf.PreInsertFuncs) != len(c.PreInsertFuncs) {
		t.Errorf("PreInsertFuncs should be overwritten. (want: %#v, got: %#v", c.PreInsertFuncs, conf.PreInsertFuncs)
	}
}
