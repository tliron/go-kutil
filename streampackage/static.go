package streampackage

//
// StaticStreamPackage
//

type StaticStreamPackage struct {
	streams []Stream

	index int
}

func NewStaticStreamPackage(streams ...Stream) *StaticStreamPackage {
	return &StaticStreamPackage{
		streams: streams,
	}
}

// StreamPackage interface
func (self *StaticStreamPackage) Next() (Stream, error) {
	if self.index < len(self.streams) {
		stream := self.streams[self.index]
		self.index++
		return stream, nil
	} else {
		return nil, nil
	}
}

// StreamPackage interface
func (self *StaticStreamPackage) Close() error {
	return nil
}
