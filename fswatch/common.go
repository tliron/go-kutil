package fswatch

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
)

var log = commonlog.GetLogger("kutil.fswatch")

type OnChangedFunc func(fileUrl *exturl.FileURL)
