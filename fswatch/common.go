package fswatch

import (
	"github.com/tliron/kutil/logging"
	urlpkg "github.com/tliron/kutil/url"
)

var log = logging.GetLogger("kutil.fswatch")

type OnChangedFunc func(fileUrl *urlpkg.FileURL)
