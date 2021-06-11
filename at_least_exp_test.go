// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"fmt"
	"time"
)

func ExampleRunAtLeast_shortTask() {
	crawler := func() error {
		// grab some web page, say it costs you 1ms
		time.Sleep(time.Millisecond)
		fmt.Println("done")
		return nil
	}

	task := RunAtLeast(10*time.Millisecond, crawler)

	task() // blocks 10ms
	task() // blocks another 10ms

	// output: done
	// done
}

func ExampleRunAtLeast_longTask() {
	crawler := func() error {
		// grab some web page, say it costs you 15ms
		time.Sleep(15 * time.Millisecond)
		fmt.Println("done")
		return nil
	}

	task := RunAtLeast(10*time.Millisecond, crawler)

	task() // blocks 15ms
	task() // blocks another 15ms

	// output: done
	// done
}

func ExampleRunSuccessAtLeast() {
	cnt := 0
	f := func() error {
		cnt++
		if cnt <= 1 {
			fmt.Println("failed")
			return errors.New("failed")
		}
		fmt.Println("done")
		return nil
	}

	task := RunSuccessAtLeast(10*time.Millisecond, f)

	task() // does not block, returns an error
	task() // blocks 10ms
	task() // blocks 10ms

	// output: failed
	// done
	// done
}

func ExampleRunFailedAtLeast() {
	cnt := 0
	f := func() error {
		cnt++
		if cnt <= 1 || cnt >= 3 {
			fmt.Println("failed")
			return errors.New("failed")
		}
		fmt.Println("done")
		return nil
	}

	task := RunFailedAtLeast(10*time.Millisecond, f)

	task() // blocks 10ms, returns an error
	task() // does not block
	task() // blocks 10ms, returns an error

	// output: failed
	// done
	// failed
}
