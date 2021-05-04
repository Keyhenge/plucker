package plucker

import (
	"errors"
	"sync"

	"github.com/rminnich/go9p"
)

const (
	EMAXMSG = 8192
)

var (
	equeue = &EQueue{}
	mutex  = sync.Mutex{} // Should this be passed in, or a global variable?
)

type EQueue struct {
	ID     int
	Buffer []byte

	//BufferLen int

	Next *EQueue
}

type Event struct {
	kbdc    int
	mouse   Mouse
	n       int
	message Message
	data    [EMAXMSG]byte
}

type Mouse struct {
	buttons int
	x       int
	y       int
	msec    uint64
}

// Should probably be renamed
func Partial(id int, event *Event, buffer []byte) error {
	var eq *EQueue
	var nmore int // Probably isn't necessary

	mutex.Lock()
	for eq = equeue; eq != nil; eq = eq.Next {
		if eq.ID == id {
			break
		}
	}
	mutex.Unlock()

	if eq == nil {
		return errors.New("Something")
	}

	//* Partial message exists for this ID
	// Message exists for this ID
	eq.Buffer = append(eq.Buffer, buffer...)
	event.message = PlumbUnpack(eq.Buffer, &nmore)

	equeue = eq.Next

	return nil
}

func AddPartial(id int, buffer []byte) {
	eq := &EQueue{
		ID:     id,
		Buffer: buffer,
	}
	mutex.Lock()
	eq.Next = equeue
	equeue = eq
	mutex.Unlock()
}

func PlumbEvent(id int, event *Event, buffer []byte) (int, error) {
	var nmore int

	if err := Partial(id, event, buffer); err != nil {
		//* No partial message already waiting for this ID
		event.message = PlumbUnpack(buffer, &nmore)
		if nmore > 0 {
			AddPartial(id, buffer)
		}
	}

	//if(event.message == nil) {
	//	return 0
	//}
	//return id

	return 0, nil
}

func EPlumb(key int, port string) error {
	var fd int

	fd = PlumbOpen(port, go9p.OREAD|go9p.OCEXEC)
	if fd < 0 {
		return errors.New("something")
	}

	return EStartFn(key, fd, PlumbEvent)
}
