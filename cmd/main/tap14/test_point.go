package tap14

import (
	"fmt"
	"io"
	"strings"
)

type TestPoint struct {
	Data string
	Directive Directive
	Label string
	Message string
	Pass bool
}

func (this *TestPoint) ReadFrom(reader io.Reader) (read int64, err error) {
	packet := newReadWritePacket()
	this.readFirstLine(packet, reader)
	if packet.Err != nil {
		goto readDone
	}

	for packet.Line() != tapBodyYAMLBlockEnd {
		scanLine(packet, reader)
		if packet.Err != nil {
			goto readDone
		}

		if isAFirstLine(packet.Line()) {
			setIncompleteError(packet)
			goto readDone
		}

		this.processLine(packet, reader)
		if packet.Err != nil {
			goto readDone
		}
	}

	readDone:
	if packet.Count > 0 && packet.Err == io.EOF {
		packet.Err = fmt.Errorf(
			"%w: %w: stream ended before finding test-point YAML body end",
			packet.Err,
			ErrIncompleteTestPoint,
		)
	}

	return packet.Count, packet.Err
}

func (this *TestPoint) WriteTo(writer io.Writer) (written int64, err error) {
	packet := newReadWritePacket()
	if this.Pass {
		packet.Wrap(fmt.Fprint(writer, tapBodyOK))
	} else {
		packet.Wrap(fmt.Fprint(writer, tapBodyNotOK))
	}

	if packet.Err != nil {
		goto writeDone
	}

	packet.Wrap(fmt.Fprint(writer, this.Label))
	if packet.Err != nil {
		goto writeDone
	}

	this.writeDirective(packet, writer)
	if packet.Err != nil {
		goto writeDone
	}

	packet.Wrap(fmt.Fprintf(writer, "%s\n", tapBodyYAMLBlockStart))
	if packet.Err != nil {
		goto writeDone
	}

	packet.Wrap(fmt.Fprintf(writer, "%s%s\n", tapBodyMessage, this.Message))
	if packet.Err != nil {
		goto writeDone
	}

	if len(this.Data) > 0 {
		packet.Wrap(fmt.Fprintf(writer, "%s%q", tapBodyData, this.Data))
		if packet.Err != nil {
			goto writeDone
		}
	}

	packet.Wrap(fmt.Fprintf(writer, "%s\n", tapBodyYAMLBlockEnd))
	if packet.Err != nil {
		goto writeDone
	}

	writeDone:
	return packet.Count, packet.Err
}

// ----- private methods -----

func (this *TestPoint) processLine(packet *readWritePacket, reader io.Reader) {
	if packet.Line() == tapBodyYAMLBlockStart {
		return
	}

	if strings.HasPrefix(packet.Line(), tapBodyMessage) {
		this.Message = strings.TrimPrefix(packet.Line(), tapBodyMessage)
		return
	}

	if strings.HasPrefix(packet.Line(), tapBodyData) {
		this.Data = strings.Trim(strings.TrimPrefix(packet.Line(), tapBodyData), "\"")
	}
}

func (this *TestPoint) readFirstLine(packet *readWritePacket, reader io.Reader) {
	var line string
	scanLine(packet, reader)

	// parse status
	if strings.HasPrefix(packet.Line(), tapBodyOK) {
		this.Pass = true
		line = strings.TrimPrefix(packet.Line(), tapBodyOK)
	} else if strings.HasPrefix(packet.Line(), tapBodyNotOK) {
		this.Pass = false
		line = strings.TrimPrefix(packet.Line(), tapBodyNotOK)
	} else if strings.HasPrefix(packet.Line(), tapBodyBail) {
		packet.Err = ErrBailOut
		return
	} else {
		packet.Err = fmt.Errorf("%w: for line \"%s\"", ErrBadStatus, packet.Line())
		return
	}

	// parse directive
	this.Directive = DirectiveNone
	if strings.HasSuffix(line, tapDirectiveSkip) {
		this.Directive = DirectiveSkip
		line = strings.TrimSuffix(line, tapDirectiveSkip)
	} else if strings.HasSuffix(line, tapDirectiveTodo) {
		this.Directive = DirectiveTodo
		line = strings.TrimSuffix(line, tapDirectiveTodo)
	}

	this.Label = line
}

func (this *TestPoint) writeDirective(packet *readWritePacket, writer io.Writer) {
	switch this.Directive {
	case DirectiveSkip:
		packet.Wrap(fmt.Fprintf(writer, "%s\n", tapDirectiveSkip))
	case DirectiveTodo:		
		packet.Wrap(fmt.Fprintf(writer, "%s\n", tapDirectiveTodo))
	default:
		packet.Wrap(fmt.Fprint(writer, "\n"))
	}
}

// ----- private functions -----

func isAFirstLine(line string) (fistLine bool) {
	return strings.HasPrefix(line, tapBodyOK) ||
		strings.HasPrefix(line, tapBodyNotOK) ||
		strings.HasPrefix(line, tapBodyBail)
}

func scanLine(packet *readWritePacket, reader io.Reader) {
	var readBuffer = []byte{0}

	packet.ResetLine()
	for {
		count := packet.Wrap(reader.Read(readBuffer))
		if packet.Err != nil && packet.Err != io.EOF {
			return
		}

		if count == 0 {
			return 
		}

		char := readBuffer[0]
		if char == '\n' {
			return
		}

		packet.AddByte(char)
	}
}

func setIncompleteError(packet *readWritePacket) {
	packet.Err = fmt.Errorf(
		"%w: starting a new record with line %q",
		ErrIncompleteTestPoint,
		packet.Line(),
	)
}