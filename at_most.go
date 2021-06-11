// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"sync"
	"time"
)

// OnceAtMost ensures every calls to f() is executed, but not too fast
//
// It guarantees:
//      - run the function as many times as you call
//      - only one call is running at the same time
//      - no more than once within the duration.
//
// Time is recorded before calling real function, which means the duration includes
// function execution time.
//
// Say you have a f() prints "test", and costs 0.1s each call
//
//    x := OnceAtMost(time.Second, f)
//    go x()
//    go x()
//    go x()
//
// You'll see a "test" every second for 3 seconds.
//
// If x() is not called in another goroutine, it acts much like RunAtLeast
//
//    x() // blocks 1s
//    x() // blocks 1s
func OnceAtMost(dur time.Duration, f func() error) func() error {
	lock := new(sync.Mutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.Lock()
		defer lock.Unlock()
		if d := time.Since(last); d <= dur {
			time.Sleep(dur - d)
		}
		last = time.Now()
		return f()
	}
}

// OnceSuccessAtMost is identical to OnceAtMost, but only successful call is ensured
//
// Say you have a f() prints "test" no matter success or failed, and costs 0.1s each call
//
//    x := OnceAtMost(time.Second, f)
//    go x() // assumes it failed
//    go x() // assumes it failed
//    go x() // assumes it succeeded
//
// You'll see:
//
//    * a "test" at 0.1s
//    * another "test" at 0.2s (0.1s after previous "test")
//    * another "test" at 1.2s (1s after previous "test")
func OnceSuccessAtMost(dur time.Duration, f func() error) func() error {
	lock := new(sync.Mutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.Lock()
		defer lock.Unlock()
		if d := time.Since(last); d <= dur {
			time.Sleep(dur - d)
		}

		now := time.Now()
		ret := f()
		if ret == nil {
			last = now
		}
		return ret
	}
}

// OnceWithin is identical to OnceAtMost, but calls within duration are ignored
//
// Say you have a f() prints "test", and costs 0.5s each call
//
//    x := OnceWithin(time.Second, f)
//    go x() // this should be executed and print "test"
//    go x() // this should be ignored
//    go x() // this should also be ignored
//    time.Sleep(time.Second)
//    go x() // this should be executed and print "test"
func OnceWithin(dur time.Duration, f func() error) func() error {
	lock := new(sync.RWMutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.RLock()
		if d := time.Since(last); d <= dur {
			lock.RUnlock()
			return nil
		}
		lock.RUnlock()

		lock.Lock()
		defer lock.Unlock()
		if d := time.Since(last); d <= dur {
			return nil
		}
		last = time.Now()
		return f()
	}
}

// OnceSuccessWithin is identical to OnceWithin, but only success call is ensured
//
// Say you have a f() prints "test", and costs 0.5s each call
//
//    x := OnceSuccessWithin(time.Second, f)
//    go x() // f executed, assumes it failed
//    go x() // f is executed, assumes it succeeded, prints "test"
//    go x() // this should be ignored
//    time.Sleep(time.Second)
//    go x() // f is executed
func OnceSuccessWithin(dur time.Duration, f func() error) func() error {
	lock := new(sync.RWMutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.RLock()
		if d := time.Since(last); d <= dur {
			lock.RUnlock()
			return nil
		}
		lock.RUnlock()

		lock.Lock()
		defer lock.Unlock()
		if d := time.Since(last); d <= dur {
			return nil
		}

		now := time.Now()
		ret := f()
		if ret == nil {
			last = now
		}
		return ret
	}
}
