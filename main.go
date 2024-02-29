package main

import (
	"github.com/onioncall/squa/common"
	"github.com/onioncall/squa/services"
	"github.com/onioncall/squa/welcome"
)

func main() {
	services.Clear()
	env := common.GetEnvironment()
	common.SetEnvironment(env.Production)

	// services.LatestMessageId = 0
	welcome.Execute()
}
