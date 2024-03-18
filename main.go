package main

import (
	"os"

	"github.com/onioncall/squa/common"
	"github.com/onioncall/squa/services"
	"github.com/onioncall/squa/welcome"
)

func main() {
	services.Clear()

	env := os.Getenv("ENV")
	if env == "Development" { 
		common.SetEnvironment("http://localhost:8040")
	} else {
		common.SetEnvironment("https://clisqua.com")
	}
	
	// set timer for five minutes before quitting to save on azure costs.
	// remove when we have a more sensible database option
	welcome.Execute()
}
