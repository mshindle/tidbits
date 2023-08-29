package leader

type Event string

type HandlerFunc func(event Event, payload any)

type EventBus struct {
	h map[Event][]HandlerFunc
}

type MessageType uint32

const (
	LeaderElected Event       = "leaderElected"
	PING          MessageType = iota + 1
	PONG
	ELECTION
	ALIVE
	ELECTED
	OK
)

func NewEventBus() *EventBus {
	return &EventBus{
		h: make(map[Event][]HandlerFunc),
	}
}

func (e *EventBus) Subscribe(event Event, handlers ...HandlerFunc) {
	e.h[event] = append(e.h[event], handlers...)
}

func (e *EventBus) Emit(event Event, payload any) {
	for _, handler := range e.h[event] {
		go handler(event, payload)
	}
}

type Message struct {
	FromPeerID string
	Rank       int
	Type       MessageType
}

func (m *Message) IsAliveMessage() bool {
	return m.Type == ALIVE
}

func (m *Message) IsPongMessage() bool {
	return m.Type == PONG
}
