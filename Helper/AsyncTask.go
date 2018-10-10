package helper

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	CanceledError       = errors.New("Operation canceled")
	AlreadyStartedError = errors.New("Operation already running")
)

type AsyncTask interface {
	// Control methods
	Start()
	Stop()
	// Query methods
	Status() Status
}

type Status struct {
	Running  bool
	Error    error
	Progress interface{}
}

type Link interface {
	Canceled() bool
	SetProgressInfo(interface{})
}

/*
type ProgressReader interface {
	Get() Progress
}
*/

func MakeAsyncTaskFunc(fn func(Link) error) AsyncTask {
	return &asyncFunc{taskFn: fn}
}

type asyncFunc struct {
	taskFn func(Link) error

	mutex         sync.Mutex
	started       bool
	cancelPending atomic.Value
	err           error
	progressInfo  atomic.Value
}

func (obj *asyncFunc) Start() {
	// Check already one running
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if obj.started {
		return //AlreadyStartedError
	}
	// Set state to running
	obj.started = true
	obj.cancelPending.Store(false)
	obj.err = nil

	go func() {
		var err error
		defer func() {
			obj.mutex.Lock()
			// Recover if panic
			if x := recover(); x != nil {
				obj.err = fmt.Errorf("run time panic: %v", x)
			} else if obj.cancelPending.Load().(bool) {
				// TODO prendre en compte si erreur retournée par la méthode ??
				obj.err = CanceledError
			} else {
				obj.err = err
			}
			// In last set state to runnable
			obj.started = false
			obj.mutex.Unlock()
		}()

		// Call the task function
		err = obj.taskFn(obj)
	}()
}

func (obj *asyncFunc) Stop() {
	// It's not important to check if running or not, because when Start, cancelPending.Store(false)
	obj.cancelPending.Store(true)
}

func (obj *asyncFunc) Status() Status {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return Status{
		Running:  obj.started,
		Error:    obj.err,
		Progress: obj.progressInfo.Load(),
	}
}

func (obj *asyncFunc) Canceled() bool {
	return obj.cancelPending.Load().(bool)
}

func (obj *asyncFunc) SetProgressInfo(data interface{}) {
	obj.progressInfo.Store(data)
}

/*
// Todo passage à State pour remplacer: ???
	- started       bool
	- cancelPending atomic.Value

type State int

const (
	None State = iota + 1
	Running
	CancelPending
	Completed
)
*/
