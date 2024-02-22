package main

import (
	"github.com/onioncall/cli-squa/cli/common"
	"github.com/onioncall/cli-squa/cli/services"
	"github.com/onioncall/cli-squa/cli/welcome"
)

func main() {
	services.Clear()
	env := common.GetEnvironment()
	common.SetEnvironment(env.Develop)

	common.LatestMessageId = 0
	welcome.Execute()
}

