// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

// Recorded creates a function that remembers how many times it is called
//
// for example:
//
//    x := Recorded(f)
//    x() // equals f(0)
//    x() // equals f(1)
//    x() // equals f(2)
func Recorded(f func(idx uint64) error) (ret func() error) {
	tries := uint64(0)

	return func() (err error) {
		err = f(tries)
		tries++
		return
	}
}

// Retry runs f() until it returns nil
//
// err is closed when f() returns nil
//
// common usacase:
//
//    ch := Retry(RunAtLeast(time.Minute, f)) // retries every minute
//    idx := 1
//    for e := range ch {
//        log.Printf("#%d attempt is failed: %v", idx, err)
//        idx++
//    }
//    log.Print("#%d attempt is successfully done", idx)
//
// WARNING: err is not buffered, so it won't execute before error in err is consumed
func Retry(f func() error) (err chan error) {
	err = make(chan error)

	go func(err chan error) {
		defer close(err)
		for {
			e := f()
			if e == nil {
				return
			}

			err <- e
		}
	}(err)

	return
}

// TriesAtMost retries f for at most n times
//
// suppose your f() costs 1s to run and always fail, TriesAtMost(10, f) will run f()
// 10 times (which costs you 10s), and err is closed at 10s
func TriesAtMost(n uint64, f func() error) (err chan error) {
	return Retry(Recorded(func(idx uint64) error {
		if idx >= n {
			return nil
		}

		return f()
	}))
}

// IgnoreErr drops all errors in ch asynchronously
//
// If you need it synchronously, just use "for range ch {}".
func IgnoreErr(ch chan error) {
	go func() {
		for range ch {
		}
	}()
}
