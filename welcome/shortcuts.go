package welcome

import (
	"github.com/onioncall/squa/create"
	"github.com/onioncall/squa/join"
)

func RouteShortcut(input string) {
	switch input {
	case ".c":
		create.Execute()
	case ".j":
		join.Execute()
	}
}