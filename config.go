package mongotest

import (
	"errors"
	"path/filepath"
)

// FixtureFormatType is decision policy of fixture data format.
type FixtureFormatType string

const (
	// FixtureFormatAuto means that fixture format is decided with file extension. (default)
	FixtureFormatAuto = FixtureFormatType("Auto")
	// FixtureFormatJSON means that fixture is written with JSON format.
	FixtureFormatJSON = FixtureFormatType("JSON")
	// FixtureFormatYAML means that fixture is written with YAML format.
	FixtureFormatYAML = FixtureFormatType("YAML")
	// fixtureFormatUnknown means that fixture is written with unknown format.
	// Not export this value. Using error instead of this value.
	fixtureFormatUnknown = FixtureFormatType("Unknown")

	defaultTimeoutSeconds = 10
)

type config struct {
	url               string
	database          string
	fixtureRootDir    string
	fixtureRootDirAbs string
	fixtureFormat     FixtureFormatType
	timeout           int
	preInsertFuncs    []PreInsertFunc
}

func defaultConfig() *config {
	return &config{
		fixtureFormat:  FixtureFormatAuto,
		timeout:        defaultTimeoutSeconds,
		preInsertFuncs: make([]PreInsertFunc, 0),
	}
}

// ConfigFunc is function for setting configuration parameter.
type ConfigFunc func(conf *config) *config

// PreInsertFunc is function for doing additional action to values.
type PreInsertFunc func(collectionName string, doc DocData) (DocData, error)

// URL returns function for setting MongoDB server url.
//   ex: mongodb://localhost:27017
func URL(url string) ConfigFunc {
	return func(conf *config) *config {
		conf.url = url
		return conf
	}
}

// Database returns function for setting database name that is connected to.
func Database(name string) ConfigFunc {
	return func(conf *config) *config {
		conf.database = name
		return conf
	}
}

// FixtureRootDir returns function for setting root directory path of fixtures.
func FixtureRootDir(dir string) ConfigFunc {
	return func(conf *config) *config {
		conf.fixtureRootDir = dir
		return conf
	}
}

// FixtureFormat returns function for setting format of fixtures.
func FixtureFormat(format FixtureFormatType) ConfigFunc {
	return func(conf *config) *config {
		conf.fixtureFormat = format
		return conf
	}
}

// Timeout returns function for setting timeout seconds.
// If given value is minus, this func does not set value to config.
// Default timeout is 10 seconds.
func Timeout(timeout int) ConfigFunc {
	return func(conf *config) *config {
		conf.timeout = timeout
		return conf
	}
}

// PreInsert returns function for setting function for doing additional action to values.
func PreInsert(fn PreInsertFunc) ConfigFunc {
	return func(conf *config) *config {
		conf.preInsertFuncs = append(conf.preInsertFuncs, fn)
		return conf
	}
}

func validateConfig() error {
	if conf.url == "" {
		return errors.New("empty URL")
	}
	if conf.database == "" {
		return errors.New("empty database name")
	}
	if conf.timeout <= 0 {
		return errors.New("invalid timeout seconds")
	}
	abs, err := filepath.Abs(conf.fixtureRootDir)
	if err != nil {
		return err
	}
	conf.fixtureRootDirAbs = abs
	return nil
}
