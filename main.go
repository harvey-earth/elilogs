package main

import (
	"github.com/harvey-earth/elilogs/cmd"
	"github.com/harvey-earth/elilogs/utils"
)

func main() {
	cmd.Configure()
	utils.InitLogger()
	cmd.Execute()
}
