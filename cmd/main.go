package main

import (
	"fmt"

	configEnv "github.com/caarlos0/env/v10"
)

type DbConfig struct {
	Username string `env:"USERNAME,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
}

func main() {
	var dbVars DbConfig
	err := configEnv.Parse(&dbVars)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", dbVars)
}
