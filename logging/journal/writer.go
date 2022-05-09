package journal

import (
	"github.com/coreos/go-systemd/journal"
	"github.com/tliron/kutil/util"
)

//
// JournalWriter
//

type JournalWriter struct{}

// io.Writer interface
func (self JournalWriter) Write(p []byte) (int, error) {
	journal.Send(util.BytesToString(p), journal.PriDebug, nil)
	return len(p), nil
}
