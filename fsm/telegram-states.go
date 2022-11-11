package fsm

// to initialize state machine + states

var FSM = NewStateMachine()

var HelloState = FSM.NewState("hello")
var ByeState = FSM.NewState("buy")
