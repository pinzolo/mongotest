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
	fixtureFormatEmpty   = FixtureFormatType("")

	defaultTimeoutSeconds = 10
)

// Config is configuration holder of mongotest module.
type Config struct {
	URL            string
	Database       string
	FixtureRootDir string
	FixtureFormat  FixtureFormatType
	Timeout        int
	PreInsertFuncs []PreInsertFunc

	fixtureRootDirAbs string
}

var conf = defaultConfig()

// Configuration returns current config.
func Configuration() Config {
	return conf
}

func defaultConfig() Config {
	return Config{
		FixtureFormat:  FixtureFormatAuto,
		Timeout:        defaultTimeoutSeconds,
		PreInsertFuncs: make([]PreInsertFunc, 0),
	}
}

// PreInsertFunc is function for doing additional action to values.
type PreInsertFunc func(collectionName string, doc DocData) (DocData, error)

func validateConfig() error {
	if conf.URL == "" {
		return errors.New("empty URL")
	}
	if conf.Database == "" {
		return errors.New("empty Database name")
	}
	if conf.Timeout <= 0 {
		return errors.New("invalid Timeout seconds")
	}
	abs, err := filepath.Abs(conf.FixtureRootDir)
	if err != nil {
		return err
	}
	conf.fixtureRootDirAbs = abs
	return nil
}

// Configure overwrite configuration by given config.
func Configure(c Config) {
	if c.URL != "" {
		conf.URL = c.URL
	}
	if c.Database != "" {
		conf.Database = c.Database
	}
	if c.FixtureRootDir != "" {
		conf.FixtureRootDir = c.FixtureRootDir
	}
	if c.FixtureFormat != fixtureFormatEmpty {
		conf.FixtureFormat = c.FixtureFormat
	}
	if c.Timeout > 0 {
		conf.Timeout = c.Timeout
	}
	if c.PreInsertFuncs != nil {
		conf.PreInsertFuncs = c.PreInsertFuncs
	}
}
