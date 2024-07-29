package access

import (
	"fmt"
	"sort"
	"testing"

	"github.com/eolinker/go-common/access"
)

func TestPrintlnRoleAccess(t *testing.T) {
	system, has := access.GetPermit("system")
	if has {
		keys := system.AccessKeys()
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("- %s\n", k)
		}

	}
	team, has := access.GetPermit("team")
	if has {
		keys := team.AccessKeys()
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("- %s\n", k)
		}
	}
}
