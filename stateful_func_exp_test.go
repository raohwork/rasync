// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"fmt"
	"time"
)

func ExampleStatefulFunc_isRunning() {
	f := func() error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("run")
		return nil
	}

	sf := NewStatefulFunc(f)
	fmt.Println(sf.IsRunning()) // false
	go sf.Run()                 // "run" at 10ms
	time.Sleep(5 * time.Millisecond)
	fmt.Println(sf.IsRunning()) // "true" at 5ms
	time.Sleep(6 * time.Millisecond)
	fmt.Println(sf.IsRunning()) // "false" as 11ms

	// output: false
	// true
	// run
	// false
}

func ExampleStatefulFunc_tryRun() {
	f := func() error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("run")
		return nil
	}

	sf := NewStatefulFunc(f)
	go sf.TryRun() // "run" at 10ms
	time.Sleep(time.Millisecond)
	go sf.Run() // "run" at 20ms
	time.Sleep(5 * time.Millisecond)
	fmt.Println(sf.TryRun()) // ErrRunning at 6ms
	time.Sleep(10 * time.Millisecond)
	fmt.Println(sf.TryRun()) // ErrRunning at 16ms
	time.Sleep(5 * time.Millisecond)
	fmt.Println(sf.TryRun()) // called at 21ms, so print "run" and "<nil>" at 31ms

	// output: StatefulFunc: function is running
	// run
	// StatefulFunc: function is running
	// run
	// run
	// <nil>
}

func ExampleStatefulFunc_lock() {
	f := func() error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("run")
		return nil
	}

	sf := NewStatefulFunc(f)

	go func() {
		sf.Run() // "run" at 10ms
		sf.Run() // blocked by Lock() for 5ms, so print "run" at 25ms
		sf.Run() // "run" at 35ms
	}()

	time.Sleep(1 * time.Millisecond)
	release := sf.Lock() // return at 10ms
	fmt.Println("lock")  // printed at 10ms, just after first sf.Run() returns
	time.Sleep(5 * time.Millisecond)
	release() // release the lock so second Run() can be executed

	time.Sleep(30 * time.Millisecond)

	// output: run
	// lock
	// run
	// run
}
