package ui

import (
	"fmt"
	"sync"
)

var DefaultWgtMgr WgtMgr

var counter struct {
	sync.RWMutex
	count int
}

type Widget interface {
	Id() string
}

type WgtMgr map[EventType]WgtInfo

type WgtInfo struct {
	Handlers map[EventType]func(Event)
	WgtRef   Widget
	Id       string
}

func NewWgtMgr() WgtMgr {
	wm := WgtMgr(make(map[EventType]WgtInfo))
	return wm

}

func (wm WgtMgr) WgtHandlersHook() func(Event) {
	return func(e Event) {
		for _, v := range wm {
			if f := findMatch(v.Handlers, e.Type); f != nil {
				f(e)
			}
		}
	}
}

func findMatch(mux map[EventType]func(Event), path EventType) func(Event) {
	f, ok := mux[path]
	if ok {
		return f
	}
	return nil
}

func isPathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		return false
	}
	n := len(pattern)
	return len(path) >= n && path[0:n] == pattern
}

func GenId() string {
	counter.Lock()
	defer counter.Unlock()

	counter.count += 1
	return fmt.Sprintf("%d", counter.count)
}


