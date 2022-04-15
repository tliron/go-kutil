package exec

import (
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("kutil.exec")

type Size struct {
	Width  uint
	Height uint
}
