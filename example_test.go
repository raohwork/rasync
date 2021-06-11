// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"fmt"
	"time"
)

func Example() {
	// our function
	x := 0
	f := func() error {
		// do some time-consume task like calling external api
		x++
		fmt.Printf("%d ", x)

		return nil
	}

	// - run the function 3 times
	// - at least 50ms between two executions
	f1 := OnceAtMost(50*time.Millisecond, f)
	f1() // 1
	f1() // 2
	f1() // 3

	// - run the function as many times as possible for 50ms
	// - at least 10ms between two executions
	//
	// WARNING: THIS IS NOT A GOOD APPROACH IF f() RUNS FAST.
	begin := time.Now()
	f2 := OnceWithin(10*time.Millisecond, f)
	for time.Since(begin) < 50*time.Millisecond {
		f2() // 4 5 6 7 8
	}

	// function always spend at least 10ms
	RunAtLeast(10*time.Millisecond, f)() // 9

	// Output: 1 2 3 4 5 6 7 8 9
}
