package tap14

import (
	"encoding/json"
	"fmt"
)

// ----- Data Option -----

func SetData(data any) TestOption {
	return func(testPoint *TestPoint) (err error) {
		if data == nil {
			testPoint.Data = "null"
			return nil
		}

		var jsonBytes []byte
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			return err
		}

		testPoint.Data = fmt.Sprintf("%q", string(jsonBytes))
		return nil
	}
}

// ----- Directive Options -----

func Skip(testPoint *TestPoint) (err error) {
	testPoint.Directive = DirectiveSkip
	return nil
}

func Todo(testPoint *TestPoint) (err error) {
	testPoint.Directive = DirectiveTodo
	return nil
}

// ----- Label Option -----

func AddLabel(label string) TestOption {
	return func(testPoint *TestPoint) (err error) {
		if len(testPoint.Label) == 0 {
			testPoint.Label = label
			return nil
		}

		testPoint.Label = fmt.Sprintf("%s - %s", testPoint.Label, label)
		return nil
	}
}

// ----- Message Option -----

func AddMessage(message string) TestOption {
	return func(testPoint *TestPoint) (err error) {
		if len(testPoint.Message) == 0 {
			testPoint.Message = fmt.Sprintf("%q", message)
			return nil
		}

		testPoint.Message = fmt.Sprintf("%s - %q", testPoint.Message, message)
		return nil
	}
}
