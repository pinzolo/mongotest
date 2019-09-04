package mongotest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pinzolo/mongotest"
)

func TestMain(m *testing.M) {
	mongotest.Configure(
		mongotest.URL("mongodb://root:password@127.0.0.1:27017"),
		mongotest.Database("mongotest"),
		mongotest.FixtureRootDir("testdata"),
		mongotest.PreInsert(mongotest.SimpleConvertTime("users", "created_at")))

	if err := mongotest.Try(); err != nil {
		fmt.Println("Cannot connect to database, please run `docker-compose up -d`")
		os.Exit(2)
	}
	code := m.Run()
	os.Exit(code)
}
