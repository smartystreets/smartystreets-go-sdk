package tap14

import "errors"

var (
	ErrBailOut = errors.New(tapBodyBail)
	ErrBadStatus = errors.New("unknown status")
	ErrIncompleteTestPoint = errors.New("incomplete test point")

	ErrBadEmitterState = errors.New("bad emitter state")
)