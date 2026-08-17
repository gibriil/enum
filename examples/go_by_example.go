package go_by_example_main

import (
	"fmt"

	"github.com/gibriil/enum"
)

type ServerState int

const (
	StateIdle ServerState = iota
	StateConnected
	StateError
	StateRetrying
)

var ServerStates = enum.DefineType(
	enum.As[ServerState]{Name: "idle", Value: StateIdle},
	enum.As[ServerState]{Name: "connected", Value: StateConnected},
	enum.As[ServerState]{Name: "error", Value: StateError},
	enum.As[ServerState]{Name: "retrying", Value: StateRetrying},
)

func (ss ServerState) String() string {
	e := enum.Of(ss)
	return e.Name()
}

func main() {
	ns := transition(enum.Of(StateIdle))
	fmt.Println(ns)

	ns2 := transition(enum.Of(ns))
	fmt.Println(ns2)
}

func transition(s enum.EnumAs[ServerState]) ServerState {
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
