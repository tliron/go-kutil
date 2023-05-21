package exec

import (
	"github.com/tliron/commonlog"
)

var log = commonlog.GetLogger("kutil.exec")

type Size struct {
	Width  uint
	Height uint
}
