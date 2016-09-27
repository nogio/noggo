package base

type Type string


const (
	TypeTriggerNone			= "none"
	TypeTriggerFinish		= "finish"
	TypeTriggerRetrigger	= "retrigger"
)

const (
	TypeErrorNone = "none"
	TypeErrorFailed = "failed"
	TypeErrorDenied = "denied"
	TypeErrorState = "state"
)