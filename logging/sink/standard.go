package sink

import (
	"bufio"
	"io"
	"log"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/util"
)

type StandardLogParseFunc func(line string) logging.Message

func NewStandardLogger(parse StandardLogParseFunc) *log.Logger {
	pipeReader, pipeWriter := io.Pipe()
	util.OnExitError(pipeReader.Close)

	go func() {
		reader := bufio.NewReader(pipeReader)
		for {
			if line, err := reader.ReadString('\n'); err == nil {
				if len(line) > 1 {
					line = line[:len(line)-1]
					if message := parse(line); message != nil {
						message.Send()
					}
				}
			} else {
				return
			}
		}
	}()

	return log.New(pipeWriter, "", 0)
}
