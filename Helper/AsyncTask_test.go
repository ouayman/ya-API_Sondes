package helper

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestAsyncTask_PanicRecover(t *testing.T) {
	// Panic recover
	panicFunc := MakeAsyncTaskFunc(func(link Link) error {
		panic("Collector: Panic testing")
		return nil
	})
	panicFunc.Start()
	time.Sleep(10 * time.Millisecond)
	status := panicFunc.Status()
	if status.Running {
		t.Error("Panic: not running attended")
	}
	if status.Error == nil {
		t.Error("Panic: error attended")
	}

	// Panic type assertion recover
	assertErrorFunc := MakeAsyncTaskFunc(func(link Link) error {
		var t interface{}
		println(t.(string))
		return nil
	})
	assertErrorFunc.Start()
	time.Sleep(10 * time.Millisecond)
	status = assertErrorFunc.Status()
	if status.Running {
		t.Error("Type assertion error: not running attended")
	}
	if status.Error == nil {
		t.Error("Type assertion error: error attended")
	}
}

func TestAsyncTask_LifeState(t *testing.T) {
	// Test status running and no error
	taskFunc := MakeAsyncTaskFunc(func(link Link) error {
		start := time.Now()
		for link.Canceled() == false {
			if time.Since(start).Seconds() > 10 {
				break
			}
			time.Sleep(time.Millisecond)
		}
		return nil
	})
	// Check initial state
	status := taskFunc.Status()
	if status.Error != nil {
		t.Error("Initial: error nil attended")
	}
	// Start and check immediatly task running
	taskFunc.Start()
	status = taskFunc.Status()
	if status.Running == false {
		t.Error("Running: running attended")
	}
	if status.Error != nil {
		t.Error("Running: no error attended")
	}
	// Test after 0.5 seconds always running
	time.Sleep(500 * time.Millisecond)
	status = taskFunc.Status()
	if status.Running == false {
		t.Error("Running: running attended")
	}
	if status.Error != nil {
		t.Error("Running: no error attended")
	}
	// Stop and check after 2 milliseconds task not running
	taskFunc.Stop()
	time.Sleep(2 * time.Millisecond)
	status = taskFunc.Status()
	if status.Running {
		t.Error("Running: not running attended")
	}
	if status.Error != CanceledError {
		t.Error("Running: canceled error attended")
	}
}

func TestAsyncTask_ReturnError(t *testing.T) {
	// Error returned
	var errBack = errors.New("Error returned by function")
	returnErrorFunc := MakeAsyncTaskFunc(func(link Link) error {
		return errBack
	})
	returnErrorFunc.Start()
	time.Sleep(10 * time.Millisecond)
	status := returnErrorFunc.Status()
	if status.Running {
		t.Error("Error returned: not running attended")
	}
	if status.Error != errBack {
		t.Error("Error returned: error attended")
	}
}

func TestAsyncTask_AccessToData(t *testing.T) {
	// Test simultaneous progress access (read and write)
	readAndWriteAccessFunc := MakeAsyncTaskFunc(func(link Link) error {
		var progress Progress
		progress.Total = 100000
		for i := 0; i < progress.Total; i++ {
			progress.Message = strconv.Itoa(i)
			progress.Current++
			link.SetProgressInfo(progress)
		}
		return nil
	})
	readAndWriteAccessFunc.Start()
	for status := readAndWriteAccessFunc.Status(); status.Running; status = readAndWriteAccessFunc.Status() {
		_, err := json.Marshal(status.Progress)
		if err != nil {
			t.Error("Read And Write Access: error: ", err.Error())
		}
	}
	// Access when finished
	status := readAndWriteAccessFunc.Status()
	_, err := json.Marshal(status.Progress)
	if err != nil {
		t.Error("Read And Write Access: error when finished: ", err.Error())
	}
}
