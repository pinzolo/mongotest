package mongotest

var DefaultTimeoutSeconds = defaultTimeoutSeconds

func Reconfigure(opts ...ConfigFunc) (reset func()) {
	orig := *conf
	if len(opts) > 0 {
		Configure(opts...)
	}
	return func() {
		conf = &orig
	}
}

func DefaultConfig() (reset func()) {
	orig := conf
	conf = defaultConfig()
	return func() {
		conf = orig
	}
}

func GetConfigURL() string {
	return conf.url
}

func GetConfigDatabase() string {
	return conf.database
}

func GetConfigFixtureRootDir() string {
	return conf.fixtureRootDir
}

func GetConfigFixtureFormat() FixtureFormatType {
	return conf.fixtureFormat
}

func GetConfigTimeout() int {
	return conf.timeout
}
