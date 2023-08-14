package streampackage

import (
	contextpkg "context"
	"io"
)

//
// StreamPackage
//

type StreamPackage interface {
	Next() (Stream, error)
	Close() error
}

//
// Stream
//

type Stream interface {
	Open(context contextpkg.Context) (io.ReadCloser, string, bool, error) // path, isExecutable
}
