package base

type Type int


const (
	TypeTriggerNone = iota
	TypeTriggerFinish
	TypeTriggerRetrigger
)

const (
	TypeErrorNone = iota
	TypeErrorFailed
	TypeErrorDenied
	TypeErrorState
)