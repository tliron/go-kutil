package exec

//
// Process
//

type Process struct {
	Stdout chan []byte // receive from this
	Stderr chan []byte // receive from this

	stdin  chan []byte // send to this
	resize chan Size   // send to this
}

func newProcess(channelSize int) Process {
	return Process{
		Stdout: make(chan []byte, channelSize),
		Stderr: make(chan []byte, channelSize),
		stdin:  make(chan []byte, channelSize),
		resize: make(chan Size, channelSize),
	}
}

// Must be called, otherwise there will be a goroutine leak!
func (self *Process) Close() {
	close(self.stdin)
	//close(self.resize)
}

func (self *Process) Stdin(p []byte) {
	if p != nil {
		self.stdin <- p
	}
}

func (self *Process) Resize(width uint, height uint) {
	self.resize <- Size{Width: width, Height: height}
}
