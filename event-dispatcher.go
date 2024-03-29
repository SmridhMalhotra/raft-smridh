package main

import (
	"reflect"
	"sync"
)

type eventDispatcher struct{
	sync.RWMutex
	source interface{}
	listeners map[string]eventListeners
}

type EventListener func(Event)

type eventListeners []EventListener

func newEventDispatcher(source interface{}) *eventDispatcher{
	return &eventDispatcher{
		source:    source,
		listeners: make(map[string]eventListeners),
	}
}

func (d *eventDispatcher) AddEventListener(typ string, listener EventListener){
	d.Lock()
	defer d.Unlock()
	d.listeners[typ] = append(d.listeners[typ],listener)
}

func (d *eventDispatcher) RemoveEventListener(typ string,listener EventListener){
	d.Lock()
	defer d.Unlock()
	ptr := reflect.ValueOf(listener).Pointer()

	listeners := d.listeners[typ]
	for i,x := range listeners{
		if reflect.ValueOf(x).Pointer() == ptr {
			d.listeners[typ] = append(listeners[:i],listeners[i+1:]...)
		}
	}
}

func (d *eventDispatcher) DispatchEvent(e Event){
	d.RLock()
	defer d.RUnlock()

	if e,ok := e.(*event);ok{
		e.source = d.source
		for _,l := range d.listeners[e.Type()]{
			l(e)
		}
	}


}
