// Package go_by_example_main demonstrates type-backed enum definitions.
package go_by_example_main

import (
	"fmt"

	"github.com/gibriil/enum"
)

// ServerState identifies a server connection state.
type ServerState int

// ServerState values registered in ServerStates.
const (
	StateIdle ServerState = iota
	StateConnected
	StateError
	StateRetrying
)

// ServerStates is the namespace for the server-state values.
var ServerStates = enum.DefineType(
	enum.As[ServerState]{Name: "idle", Value: StateIdle},
	enum.As[ServerState]{Name: "connected", Value: StateConnected},
	enum.As[ServerState]{Name: "error", Value: StateError},
	enum.As[ServerState]{Name: "retrying", Value: StateRetrying},
)

// String returns the registered name of the server state.
func (ss ServerState) String() string {
	e := enum.Of(ss)
	return e.Name()
}

// main demonstrates transitions between server states.
func main() {
	ns := transition(enum.Of(StateIdle))
	fmt.Println(ns)

	ns2 := transition(enum.Of(ns))
	fmt.Println(ns2)
}

// transition returns the next state for a registered server-state member.
func transition(s enum.Member[ServerState]) ServerState {
	switch s.Raw() {
	case StateIdle:
		return StateConnected
	case StateConnected, StateRetrying:
		return StateIdle
	case StateError:
		return StateError
	default:
		panic(fmt.Errorf("unknown state: %s", s))
	}
}
