package tap14

import (
	"fmt"
	"io"
)

type Emitter struct {
	started bool
	stream bool
	testPoints []TestPoint
	writer io.Writer
}

func (this *Emitter) CompleteTestGroup() (err error) {
	if !this.started {
		return fmt.Errorf(
			"%w: no test group has been started yet, common pattern is to call CompleteTestGroup from a defer block",
			ErrBadEmitterState,
		)
	}

	_, err = fmt.Fprintf(this.writer, "%s%d\n", tapPlanStart, len(this.testPoints))
	if err != nil {
		return err
	}

	if !this.stream {
		for _, testPoint := range this.testPoints {
			_, err = testPoint.WriteTo(this.writer)
			if err != nil {
				return err
			}
		}
	}

	this.testPoints = this.testPoints[:0]
	this.started = false
	return nil
}

func (this *Emitter) StartTestGroupBatched(writer io.Writer) (err error) {
	return this.start(writer, false)
}

func (this *Emitter) StartTestGroupStreamed(writer io.Writer) (err error) {
	return this.start(writer, true)
}

func (this *Emitter) TestBool(pass bool, testOptions ...TestOption) (err error) {
	var testPoint *TestPoint
	testPoint, err = generateTestPoint(pass, testOptions...)
	if err != nil {
		return err
	}

	return this.processTestPoint(testPoint)
}

func (this *Emitter) TestError(errIn error, testOptions ...TestOption) (err error) {
	pass := errIn == nil
	var testPoint *TestPoint
	testPoint, err = generateTestPoint(pass, testOptions...)
	if err != nil {
		return err
	}

	if !pass {
		errMessageOption := AddMessage(errIn.Error())
		err = errMessageOption(testPoint)
		if err != nil {
			return err
		}
	}

	return this.processTestPoint(testPoint)
}

// ----- private methods -----

func (this *Emitter) processTestPoint(testPoint *TestPoint) (err error) {
	this.testPoints = append(this.testPoints, *testPoint)

	if this.stream {
		_, err = testPoint.WriteTo(this.writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Emitter) start(writer io.Writer, stream bool) (err error) {
	if this.started {
		return fmt.Errorf(
			"%w: a test group has already been started, finish it first",
			ErrBadEmitterState,
		)
	}

	if writer == nil {
		return fmt.Errorf(
			"%w: writer is nil, use a non-nil writer",
			ErrBadEmitterState,
		)
	}

	this.started = true
	this.stream = stream
	this.writer = writer
	if len(this.testPoints) > 0 {
		this.testPoints = this.testPoints[:0]
	}

	_, err = fmt.Fprint(writer, tapVersion)
	if err != nil {
		return err
	}

	return nil
}

// ----- test functions -----

func generateTestPoint(pass bool, testOptions ...TestOption) (testPoint *TestPoint, err error) {
	testPoint = &TestPoint{
		Pass: pass,
		Directive: DirectiveNone,
	}

	for _, option := range testOptions {
		err = option(testPoint)
		if err != nil {
			return nil, err
		}
	}

	return testPoint, nil
}
