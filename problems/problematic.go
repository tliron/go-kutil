package problems

import "github.com/tliron/kutil/terminal"

//
// Problematic
//

type Problematic interface {
	Problem(*terminal.Stylist) (string, string, string, int, int)
}
