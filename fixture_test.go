package mongotest_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pinzolo/mongotest"
)

func newDiff(key string, v1, v2 interface{}) string {
	return fmt.Sprintf("%s: %v(%T) <-> %v(%T)", key, v1, v1, v2, v2)
}

const noVal = "<no value>"

var createdAt = time.Date(2019, 1, 2, 12, 34, 56, 0, time.UTC)
var createdAtPrimitive = primitive.NewDateTimeFromTime(createdAt)

func diffMap(m1 map[string]interface{}, m2 map[string]interface{}) []string {
	diffs := make([]string, 0, len(m1))
	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			if !reflect.DeepEqual(v1, v2) {
				diffs = append(diffs, newDiff(k1, v1, v2))
			}
		} else {
			diffs = append(diffs, newDiff(k1, v1, noVal))
		}
	}
	for k2, v2 := range m2 {
		if _, ok := m1[k2]; !ok {
			diffs = append(diffs, newDiff(k2, noVal, v2))
		}
	}
	return diffs
}

func TestUseFixture(t *testing.T) {
	err := mongotest.UseFixture("admin_users")
	if err != nil {
		t.Error(err)
	}
	cnt, err := mongotest.CountInt("users")
	if err != nil {
		t.Error(err)
	}
	if cnt != 2 {
		t.Errorf("saved user count is invalid (want: %d, got: %d)", cnt, 2)
	}

	saved, err := mongotest.Find("users", "admin1")
	if err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{
		"_id":        "admin1",
		"name":       "admin user1",
		"email":      "admin1@example.com",
		"admin":      true,
		"company":    "foo",
		"age":        int32(30),
		"note":       "abc",
		"created_at": createdAtPrimitive,
	}
	if !reflect.DeepEqual(saved, want) {
		t.Errorf("saved user is invalid. (%v)", diffMap(saved, want))
	}

	cnt, err = mongotest.CountInt("companies")
	if err != nil {
		t.Error(err)
	}
	if cnt != 2 {
		t.Errorf("saved user count is invalid (want: %d, got: %d)", cnt, 2)
	}
}

func TestUseFixtureWithDataSetMerge(t *testing.T) {
	err := mongotest.UseFixture("admin_users", "foo_users")
	if err != nil {
		t.Error(err)
	}
	cnt, err := mongotest.CountInt("users")
	if err != nil {
		t.Error(err)
	}
	if cnt != 3 {
		t.Errorf("saved user count is invalid (want: %d, got: %d)", 3, cnt)
	}

	saved, err := mongotest.Find("users", "admin1")
	if err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{
		"_id":        "admin1",
		"name":       "admin user1",
		"email":      "admin1@example.com",
		"admin":      true,
		"company":    "foo",
		"age":        float64(30),
		"note":       "xyz",
		"created_at": createdAtPrimitive,
	}
	if !reflect.DeepEqual(saved, want) {
		t.Errorf("saved user is invalid. (%v)", diffMap(saved, want))
	}

	cnt, err = mongotest.CountInt("companies")
	if err != nil {
		t.Error(err)
	}
	if cnt != 2 {
		t.Errorf("saved user count is invalid (want: %d, got: %d)", cnt, 2)
	}
}

func commonSuccessCheck(t *testing.T, age interface{}) {
	t.Helper()
}

func TestUseFixtureJSONFormat(t *testing.T) {
	defer mongotest.Reconfigure(mongotest.FixtureFormat(mongotest.FixtureFormatJSON))()
	err := mongotest.UseFixture("json/admin_users")
	if err != nil {
		t.Error(err)
	}
	cnt, err := mongotest.CountInt("users")
	if err != nil {
		t.Error(err)
	}
	if cnt != 2 {
		t.Errorf("saved user count is invalid (want: %d, got: %d)", cnt, 2)
	}

	saved, err := mongotest.Find("users", "admin1")
	if err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{
		"_id":        "admin1",
		"name":       "admin user1",
		"email":      "admin1@example.com",
		"admin":      true,
		"company":    "foo",
		"age":        float64(30),
		"note":       "abc",
		"created_at": createdAtPrimitive,
	}
	if !reflect.DeepEqual(saved, want) {
		t.Errorf("saved user is invalid. (%v)", diffMap(saved, want))
	}

	err = mongotest.UseFixture("yaml/admin_users")
	if err == nil {
		t.Error("should error when invalid fixture format load")
	}
}
