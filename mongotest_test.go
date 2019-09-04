package mongotest_test

import (
	"testing"

	"github.com/pinzolo/mongotest"
)

func TestTry(t *testing.T) {
	if err := mongotest.Try(); err != nil {
		t.Error(err)
	}
}

func TestTryWithEmptyURL(t *testing.T) {
	defer mongotest.DefaultConfig()()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when Database URL is empty")
	}
}

func TestTryWithEmptyDatabase(t *testing.T) {
	defer mongotest.DefaultConfig()()
	defer mongotest.Reconfigure(mongotest.Config{URL: "mongodb://root:password@localhost:27017"})()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when Database name is empty")
	}
}

func TestTryWithInvalidURL(t *testing.T) {
	defer mongotest.Reconfigure(mongotest.Config{URL: "mongodb://root:password@localhost:2701"})()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when Database URL is invalid")
	}
}
