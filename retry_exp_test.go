// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"fmt"
	"time"
)

func ExampleRecorded() {
	f := func(idx uint64) error {
		fmt.Println(idx)
		return nil
	}

	x := Recorded(f)
	x() // equals f(0)
	x() // equals f(1)
	x() // equals f(2)

	// output: 0
	// 1
	// 2
}

func ExampleRetry() {
	f := Recorded(func(idx uint64) error {
		fmt.Println(idx)
		if idx <= 1 {
			return errors.New("error")
		}
		return nil
	})

	errs := make([]error, 0, 2)
	for e := range Retry(f) {
		errs = append(errs, e)
	}

	for _, e := range errs {
		fmt.Println(e)
	}

	//output: 0
	// 1
	// 2
	// error
	// error
}

func ExampleIgnoreErr() {
	f := Recorded(func(idx uint64) error {
		fmt.Println(idx)
		if idx <= 1 {
			// handles error here
			return errors.New("error")
		}
		return nil
	})

	// errors are handled in f(), no need to process again
	IgnoreErr(Retry(f))

	// wait til goroutines exit
	time.Sleep(10 * time.Millisecond)

	//output: 0
	// 1
	// 2
}
