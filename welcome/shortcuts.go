package welcome

import (
	"github.com/onioncall/cli-squa/cli/create"
	"github.com/onioncall/cli-squa/cli/join"
)

func RouteShortcut(input string) {
	switch input {
	case ".c":
		create.Execute()
	case ".j":
		join.Execute()
	}
}