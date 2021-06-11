// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"fmt"
	"time"
)

func ExampleAllErr() {
	// crawler 1: banned at 3rd try
	crawl1 := Recorded(func(idx uint64) error {
		time.Sleep(10 * time.Millisecond)
		if idx > 1 {
			fmt.Println("banned1")
			return errors.New("banned")
		}
		fmt.Println("done1")
		return nil
	})
	// crawler 1: banned at 9th try
	crawl2 := Recorded(func(idx uint64) error {
		time.Sleep(12 * time.Millisecond)
		if idx > 7 {
			fmt.Println("banned2")
			return errors.New("banned")
		}
		fmt.Println("done2")
		return nil
	})

	ctrl := AllErr(
		InfiniteLoop(crawl2),
		InfiniteLoop(crawl1),
	)
	defer ctrl.Cancel()

	fmt.Println(<-ctrl.Err)

	// will get context canceled error from crawl2
	for err := range ctrl.Err {
		fmt.Println(err)
	}

	// output: done1
	// done2
	// done1
	// done2
	// banned1
	// banned
	// done2
	// context canceled
}

func ExampleAnyErr() {
	// crawler 1: banned at 3rd try
	crawl1 := Recorded(func(idx uint64) error {
		time.Sleep(10 * time.Millisecond)
		if idx > 1 {
			fmt.Println("banned1")
			return errors.New("banned")
		}
		fmt.Println("done1")
		return nil
	})
	// crawler 1: banned at 9th try
	crawl2 := Recorded(func(idx uint64) error {
		time.Sleep(12 * time.Millisecond)
		if idx > 7 {
			fmt.Println("banned2")
			return errors.New("banned")
		}
		fmt.Println("done2")
		return nil
	})

	ctrl := AnyErr(
		InfiniteLoop(crawl2),
		InfiniteLoop(crawl1),
	)
	defer ctrl.Cancel()

	fmt.Println(<-ctrl.Err)

	// will not get anything
	for err := range ctrl.Err {
		fmt.Println(err)
	}

	// output: done1
	// done2
	// done1
	// done2
	// banned1
	// banned
	// done2
}
