package main

type EventQueue struct {
	events []Event
}

func (qeue *EventQueue) HasNext() bool {
	return len(qeue.events) > 0
}

//NextEvent returns the next event in the qeue, a nil value denotes a empty qeue
func (queue *EventQueue) NextEvent() Event {
	if !queue.HasNext() {
		return nil
	}
	currentEvt := queue.events[0]
	queue.events = queue.events[1:]
	return currentEvt
}

//AddEvent inserts a new event into the end of the queue
func (queue *EventQueue) AddEvent(evt Event) {
	queue.events = append(queue.events, evt)
}
