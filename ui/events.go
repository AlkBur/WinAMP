package ui

import (
	"sync"
	"time"
)

// Event type. See Event.Type field.
const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventError
	EventInterrupt
	EventRaw
	EventTimer
	EventCustom
	EventNone
)

var DefaultEvtStream = NewEvtStream()
var usrEvtCh = make(chan Event)

type Event struct {
	Type   EventType // one of Event* constants
	Mod    Modifier  // one of Mod* constants or 0
	Key    Key       // one of Key* constants, invalid if 'Ch' is not 0
	Ch     rune      // a unicode character
	Width  int       // width of the screen
	Height int       // height of the screen
	Err    error     // error in case if input failed
	MouseX int       // x coord of mouse
	MouseY int       // y coord of mouse
	N      int       // number of bytes written when getting a raw event
	Data interface{}
}

type EvtStream struct {
	sync.RWMutex
	//srcMap      []chan Event
	stream      chan Event
	wg          sync.WaitGroup
	sigStopLoop chan Event
	Handlers    map[EventType][]func(Event)
	hook        func(Event)
}

type EvtTimer struct {
	Duration time.Duration
	Count    uint64
}

func NewEvtStream() *EvtStream {
	return &EvtStream{
		//srcMap:      make(map[string]chan Event),
		stream:      make(chan Event),
		Handlers:    make(map[EventType][]func(Event)),
		sigStopLoop: make(chan Event),
	}
}

func (es *EvtStream) Init() {
	es.Merge(es.sigStopLoop)
	go func() {
		es.wg.Wait()
		close(es.stream)
	}()
}

func (es *EvtStream) Merge(ec chan Event) {
	es.Lock()
	defer es.Unlock()

	es.wg.Add(1)
	//es.srcMap[name] = ec

	go func(a chan Event) {
		for n := range a {
			//n.From = name
			es.stream <- n
		}
		es.wg.Done()
	}(ec)
}

// Wait for an event and return it. This is a blocking function call.
func PollEvent() Event {
	select {
	case ev := <-input_comm:
		return ev
	case <-interrupt_comm:
		return Event{Type: EventInterrupt}
	}
}

func NewSysEvtCh() chan Event {
	ch := make(chan Event)
	go func(ch chan Event) {
		for {
			ch <- PollEvent()
		}
	}(ch)
	return ch
}

func NewTimerCh(du time.Duration) chan Event {
	t := make(chan Event)

	go func(a chan Event) {
		n := uint64(0)
		for {
			n++
			time.Sleep(du)
			e := Event{}
			e.Type = EventTimer
			e.Data = EvtTimer{
				Duration: du,
				Count:    n,
			}
			t <- e

		}
	}(t)
	return t
}

func (es *EvtStream) Handle(path EventType, handler func(Event)) {
	evs, ok := es.Handlers[path]
	if !ok {
		evs = make([]func(Event),0, 1)

	}
	evs = append(evs, handler)
	es.Handlers[path] = evs
}

func (es *EvtStream) Hook(f func(Event)) {
	es.hook = f
}

func (es *EvtStream) StopLoop() {
	go func() {
		e := Event{
			Type: EventInterrupt,
		}
		es.sigStopLoop <- e
	}()
}

func StopLoop() {
	DefaultEvtStream.StopLoop()
}

func Handle(path EventType, handler func(Event)) {
	DefaultEvtStream.Handle(path, handler)
}

func Loop() {
	DefaultEvtStream.Loop()
}

func (es *EvtStream) Loop() {
	for e := range es.stream {
		switch e.Type {
		case EventInterrupt:
			return
		}
		func(a Event) {
			es.RLock()
			defer es.RUnlock()
			fs, ok := es.Handlers[e.Type]
			if ok {
				for _, f := range fs {
					f(a)
				}
			}
		}(e)
		if es.hook != nil {
			es.hook(e)
		}
	}
}

func SendCustomEvt(data interface{}) {
	e := Event{}
	e.Type = EventCustom
	e.Data = data
	usrEvtCh <- e
}
