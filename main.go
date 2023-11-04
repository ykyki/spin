package main

import "os"

func main() {
	cli := NewCLIStd()
	os.Exit(cli.Run(os.Args))
}
