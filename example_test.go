package mongotest_test

import (
	"fmt"

	"github.com/pinzolo/mongotest"
)

func Example() {
	mongotest.Configure(mongotest.Config{
		URL:            "mongodb://root:password@127.0.0.1:27017",
		Database:       "mongotest",
		FixtureRootDir: "testdata",
		FixtureFormat:  mongotest.FixtureFormatJSON,
		PreInsertFuncs: []mongotest.PreInsertFunc{
			mongotest.SimpleConvertTime("users", "created_at"),
		},
	})

	// 1. Read testdata/json/admin_users.json and testdata/json/foo_users.json
	// 2. Merge read data
	// 3. Drop collection and insert read data
	err := mongotest.UseFixture("json/admin_users", "json/foo_users")
	if err != nil {
		panic(err)
	}

	// Count is helper function.
	// mongotest has some useful helper functions.
	n, err := mongotest.Count("users")
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	// Output: 3
}
