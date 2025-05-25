package app

var flags = &FlagsConfig{}

type FlagsConfig struct {
	Migration bool
	RunHTTP   bool
}
