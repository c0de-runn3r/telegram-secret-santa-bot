package fsm

import (
	"log"
)

type State struct {
	Name string
}

type StateMachine struct {
	CurrentState   State
	ExistingStates []State
}

type UserStateMachines struct {
	StateMachines map[string]*StateMachine
}

func NewStateMachine() *StateMachine {
	log.Print("New FSM created")
	return &StateMachine{}
}

func (sm *StateMachine) NewState(StateName string) *State {
	state := State{Name: StateName}
	sm.ExistingStates = append(sm.ExistingStates, state)
	log.Printf("New State '%v' created", StateName)
	return &state
}

func (sm *StateMachine) SetState(state State) {
	sm.CurrentState = state
	log.Printf("State '%v' is set", state)
}

func FindOrCreateUsersFSM(username string) *StateMachine {
	for k, v := range UserFSMs {
		if k == username {
			return v
		}
	}
	newFSM := NewStateMachine()
	UserFSMs[username] = newFSM
	return newFSM
}
