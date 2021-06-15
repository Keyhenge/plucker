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
func partial(id int, event *Event, buffer []byte) error {
	var eq *EQueue
	var nmore int // Probably isn't necessary
	var err error

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
	event.message, err = unpackMessage(eq.Buffer, &nmore)
	if err != nil {
		return err
	}

	equeue = eq.Next

	return nil
}

func addPartial(id int, buffer []byte) {
	eq := &EQueue{
		ID:     id,
		Buffer: buffer,
	}
	mutex.Lock()
	eq.Next = equeue
	equeue = eq
	mutex.Unlock()
}

func plumbEvent(id int, event *Event, buffer []byte) (int, error) {
	var nmore int

	if err := partial(id, event, buffer); err != nil { //probably wrong
		//* No partial message already waiting for this ID
		event.message, err = unpackMessage(buffer, &nmore)
		if err != nil {
			return -1, err //probably wrong
		}
		if nmore > 0 {
			addPartial(id, buffer)
		}
	}

	//if(event.message == nil) {
	//	return 0
	//}
	//return id

	return 0, nil
}

func ePlumb(key int, port string) error {
	fd, err := open(port, go9p.OREAD|go9p.OCEXEC)
	if err != nil {
		return err
	}
	if fd < 0 {
		return errors.New("something")
	}

	return EStartFn(key, fd, PlumbEvent)
}
