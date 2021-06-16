// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"sync"
)

var ErrRunning = errors.New("StatefulFunc: function is running")

// StatefulFunc is a function that possible to get current running state
type StatefulFunc interface {
	IsRunning() bool
	// try to run the function now, return ErrRunning immediately if it is running
	TryRun() (err error)
	// run this function, blocks until ran
	Run() (err error)
	// blocks until it's free to run, and hold the lock to prevent others running
	// calling Run() without release the lock cause deadlock! use with care.
	//
	// It's safe to call release multiple times, only first time is executed.
	Lock() (release func())
}

type statefulFunc struct {
	token chan *struct{}
	f     func() error
}

func (f *statefulFunc) IsRunning() (yes bool) {
	select {
	case x := <-f.token:
		f.token <- x
		return false
	default:
		return true
	}
}

func (f *statefulFunc) TryRun() (err error) {
	select {
	case x := <-f.token:
		err = f.f()
		f.token <- x
	default:
		err = ErrRunning
	}

	return
}

func (f *statefulFunc) Run() (err error) {
	x := <-f.token
	err = f.f()
	f.token <- x
	return
}

func (f *statefulFunc) Lock() (release func()) {
	x := <-f.token
	once := &sync.Once{}

	return func() {
		once.Do(func() {
			f.token <- x
		})
	}
}

// NewStatefulFunc creates a new StatefulFunc
func NewStatefulFunc(f func() error) (ret StatefulFunc) {
	x := &statefulFunc{
		token: make(chan *struct{}, 1),
		f:     f,
	}
	x.token <- nil

	return x
}
