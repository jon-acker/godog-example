package main

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"os"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindCommandLineFlags("godogt.", &opts)
}
