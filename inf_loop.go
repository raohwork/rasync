// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"context"
	"reflect"
)

// InfiniteLoopControl is the control structure of InfiniteLoop()
//
// Err is closed by InfiniteLoop()
type InfiniteLoopControl struct {
	Cancel context.CancelFunc
	Err    chan error
}

// AllErr merges several InfiniteLoopControl into one, and returns all error.
//
// Cancel functions are called either
//
//   - an error occurs
//   - an infinite loop exits
//
// All errors are returned. The order is determined by reflect.Select()
//
// It is possible to use it as sort of sync.WaitGroup
//
//   ctrl := AllErr(
//       InfiniteLoop(crawlsSite),
//       InfiniteLoop(crawlsAnotherSite),
//   )
//   defer ctrl.Cancel()
//
//   if err := <- ctrl.Err; err != nil {
//       log.Print("an error occurred: ", err)
//   }
//
//   // gracefully shutdown: waits all tasks to stop
//   for range ctrl.Err {
//   }
//   log.Print("all tasks are stopped, program exits")
func AllErr(ctrls ...InfiniteLoopControl) (ret InfiniteLoopControl) {
	ret = InfiniteLoopControl{
		Cancel: func() {
			for _, c := range ctrls {
				c.Cancel()
			}
		},
		Err: make(chan error),
	}

	cases := make([]reflect.SelectCase, len(ctrls))
	for idx, c := range ctrls {
		cases[idx] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c.Err),
		}
	}
	go func(ret InfiniteLoopControl, ctrls []InfiniteLoopControl, cases []reflect.SelectCase) {
		defer close(ret.Err)
		defer func() {
			for len(cases) > 0 {
				idx, val, ok := reflect.Select(cases)
				if !ok {
					// channel is closed
					cases = append(cases[:idx], cases[idx+1:]...)
					continue
				}
				ret.Err <- val.Interface().(error)
			}
		}()
		defer ret.Cancel()

		idx, val, ok := reflect.Select(cases)

		if !ok {
			// channel is closed
			cases = append(cases[:idx], cases[idx+1:]...)
			return
		}
		ret.Err <- val.Interface().(error)

	}(ret, ctrls, cases)

	return
}

// AnyErr merges several InfiniteLoopControl into one.
//
// Only the first error (could be nil if an error channel in ctrls is closed) is
// returned, later ARE DISCARDED IN BACKGROUND.
func AnyErr(ctrls ...InfiniteLoopControl) (ret InfiniteLoopControl) {
	ret = InfiniteLoopControl{
		Cancel: func() {
			for _, c := range ctrls {
				c.Cancel()
			}
		},
		Err: make(chan error),
	}

	cases := make([]reflect.SelectCase, len(ctrls))
	for idx, c := range ctrls {
		cases[idx] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c.Err),
		}
	}
	go func(ret InfiniteLoopControl, ctrls []InfiniteLoopControl) {
		defer close(ret.Err)
		// drop unread errors
		defer func() {
			for _, c := range ctrls {
				for range c.Err {
				}
			}
		}()
		defer ret.Cancel()

		_, val, ok := reflect.Select(cases)

		if !ok {
			// channel is closed
			return
		}
		ret.Err <- val.Interface().(error)

	}(ret, ctrls)

	return
}

// InfiniteLoop loops your function and capable to cancel-on-demand
//
// There are few things you should take care of:
//
//    - It will not interrupt current loop.
//    - It will not wait any second between tasks.
//
// Common usecase is InfiniteLoop(RunAtLeast(someDuration, task))
func InfiniteLoop(task func() error) (ret InfiniteLoopControl) {
	ctx, cancel := context.WithCancel(context.Background())
	err := make(chan error)
	go doInfiniteLooping(ctx, err, task)

	return InfiniteLoopControl{
		Cancel: cancel,
		Err:    err,
	}
}

func doInfiniteLooping(ctx context.Context, errchan chan error, task func() error) {
	var err error
	for err == nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			err = task()
		}
	}

	errchan <- err
	close(errchan)
}
