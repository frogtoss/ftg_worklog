package main

import (
	"fmt"
)

const VersionMajor = 1
const VersionMinor = 0

func (i *CLIVersionCmd) Run(cli *CLI) error {
	fmt.Printf("%d.%d\n", VersionMajor, VersionMinor)

	return nil
}
