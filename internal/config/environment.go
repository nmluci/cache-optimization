package config

type Environment string

const (
	EnvironmentDev  = Environment("dev")
	EnvironmentProd = Environment("prod")
)
