package tap14

type readWritePacket struct {
	Count int64
	Err error
	
	dirtyLine bool
	line string
	lineBuffer []byte
}

func newReadWritePacket() *readWritePacket {
	const bufferSize = 128
	return &readWritePacket{
		lineBuffer: make([]byte, 0, bufferSize),
		dirtyLine: true,
	}
}

func (this *readWritePacket) AddByte(char byte) {
	this.lineBuffer = append(this.lineBuffer, char)
	this.dirtyLine = true
}

func (this *readWritePacket) AddCount(count int) {
	this.Count += int64(count)
}

func (this *readWritePacket) Line() (line string) {
	if this.dirtyLine {
		this.line = string(this.lineBuffer)
		this.dirtyLine = false
	}

	return this.line
}

func (this *readWritePacket) ResetLine() {
	this.line = ""
	this.lineBuffer = this.lineBuffer[:0]
	this.dirtyLine = false
}

func (this *readWritePacket) Wrap(count int, err error) (passthrough int) {
	this.AddCount(count)
	this.Err = err

	return count
}
