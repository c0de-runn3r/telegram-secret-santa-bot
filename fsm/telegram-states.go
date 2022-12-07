package fsm

// to initialize state machine + states

var UserFSMs = make(map[string]*StateMachine)

var FSM = NewStateMachine()

var (
	ActionState              = FSM.NewState("action")
	NewGameNameState         = FSM.NewState("newGameName")
	ConnectExistingGameState = FSM.NewState("connectExistingName")
	UpdateWishesState        = FSM.NewState("updateWishes")
	MyGamesSate              = FSM.NewState("myGames")
	BudgetState              = FSM.NewState("budget")
)
