package problems

import "github.com/tliron/go-kutil/terminal"

//
// Problematic
//

type Problematic interface {
	Problem(*terminal.Stylist) (string, string, string, int, int)
}
