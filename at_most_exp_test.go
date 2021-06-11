// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"fmt"
	"time"
)

func ExampleOnceAtMost_sync() {
	f := func() error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("done")
		return nil
	}
	g := OnceAtMost(50*time.Millisecond, f)

	f() // you'll see message at 10ms
	f() // message at 20ms
	g() // message at 30ms
	g() // blocks 40ms then run f(), so message is printed at 80ms

	// output: done
	// done
	// done
	// done
}

func ExampleOnceAtMost_async() {
	f := func() error {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("done")
		return nil
	}
	g := OnceAtMost(50*time.Millisecond, f)

	go f() // you'll see message at 10ms
	go f() // another message at 10ms
	go g() // another message at 10ms
	go g() // blocks 40ms then run f(), so message is printed at 60ms

	// wait til goroutines return
	time.Sleep(100 * time.Millisecond)

	// output: done
	// done
	// done
	// done
}

func ExampleOnceWithin() {
	f := func() error {
		fmt.Println("test")
		return nil
	}

	x := OnceWithin(50*time.Millisecond, f)
	go x() // this should be executed and print "test"
	go x() // this should be ignored
	go x() // this should also be ignored
	time.Sleep(50 * time.Millisecond)
	go x() // this should be executed and print "test"

	// wait til goroutines return
	time.Sleep(50 * time.Millisecond)

	// output: test
	// test
}

func ExampleOnceSuccessWithin() {
	cnt := 0
	f := func() error {
		cnt++
		if cnt == 1 {
			fmt.Println("failed")
			return errors.New("error")
		}
		fmt.Println("test")
		return nil
	}

	x := OnceSuccessWithin(50*time.Millisecond, f)
	go x() // f executed, assumes it failed
	go x() // f is executed, assumes it succeeded, prints "test"
	go x() // this should be ignored
	time.Sleep(50 * time.Millisecond)
	go x() // f is executed

	// wait til goroutines return
	time.Sleep(50 * time.Millisecond)

	// output: failed
	// test
	// test
}
