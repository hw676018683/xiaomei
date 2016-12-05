package oam

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/cli/cli"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Shell() {
	for _, addr := range cli.MatchedServerAddrs() {
		address := config.Data.DeployUser + `@` + addr
		cmd.Run(cmd.O{Panic: true}, `ssh`, `-t`, address, `cd `+config.Data.DeployPath+`; bash`)
		return
	}
}