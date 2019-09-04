package mongotest

var DefaultTimeoutSeconds = defaultTimeoutSeconds

func Reconfigure(c Config) (reset func()) {
	orig := conf
	Configure(c)
	return func() {
		conf = orig
	}
}

func DefaultConfig() (reset func()) {
	orig := conf
	conf = defaultConfig()
	return func() {
		conf = orig
	}
}
