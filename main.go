package main

import (
	"fmt"

	"github.com/ach1000/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	err = cfg.SetUser("andy")
	if err != nil {
		panic(err)
	}

	updatedCfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", updatedCfg)
}
