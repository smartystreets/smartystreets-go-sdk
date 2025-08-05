// tap14 package provides emitters and consumers for the TAP-14 (Test Anything
// Protocol v14). The specification can be found here:
// https://testanything.org/tap-version-14-specification.html
package tap14

const (
	DirectiveNone Directive = iota
	DirectiveSkip
	DirectiveTodo
)

const (
	tapDirectiveSkip = " # SKIP "
	tapDirectiveTodo = " # TODO "
	
	tapBodyBail = "Bail out! "
	tapBodyData = "  data: "
	tapBodyMessage = "  message: "
	tapBodyNotOK = "\033[91mnot ok\033[0m - "
	tapBodyOK = "\033[92mok\033[0m - "
	tapBodyYAMLBlockEnd = "  ..."
	tapBodyYAMLBlockStart = "  ---"
	
	tapPlanStart = "1.."
	
	tapVersion = "TAP version 14"
)

type Directive int

type TestOption func(restPoint *TestPoint) (err error)
