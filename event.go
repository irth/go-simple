package simple

type EventType int

const (
	UnknownEventType  EventType = 0
	InputEventType    EventType = 1
	SelectedEventType EventType = 2
	RenderedEventType EventType = 3
)

func (e EventType) String() string {
	switch e {
	case InputEventType:
		return "input"
	case SelectedEventType:
		return "selected"
	case RenderedEventType:
		return "rendered"
	default:
		return "unknown"
	}
}

type Event interface {
	Type() EventType
}

type UnknownEvent struct {
	Raw string
}

func (UnknownEvent) Type() EventType { return UnknownEventType }

type InputEvent struct {
	ID    string
	Value string
}

func (InputEvent) Type() EventType { return InputEventType }

type SelectedEvent struct {
	ID string
}

func (SelectedEvent) Type() EventType { return SelectedEventType }

type RenderedEvent struct{}

func (RenderedEvent) Type() EventType { return RenderedEventType }
