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
	defer mongotest.Reconfigure(mongotest.URL(""))()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when database URL is empty")
	}
}

func TestTryWithEmptyDatabase(t *testing.T) {
	defer mongotest.Reconfigure(mongotest.Database(""))()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when database name is empty")
	}
}

func TestTryWithMinusTimout(t *testing.T) {
	defer mongotest.Reconfigure(mongotest.Timeout(-1))()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when timeout seconds is invalid")
	}
}

func TestTryWithInvalidURL(t *testing.T) {
	defer mongotest.Reconfigure(mongotest.URL("mongodb://root:password@localhost:2701"))()
	if err := mongotest.Try(); err == nil {
		t.Error("Try should return error when database URL is invalid")
	}
}
