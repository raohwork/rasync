// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import "time"

// RunAtLeast ensures the execution time is greater than the duration
//
// It will blocked until dur is reached and f() is returned.
func RunAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Since(begin); d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}

// RunSuccessAtLeast is identical to RunAtLeast, but only successful call is ensured.
//
// Say your f() needs 0.1s no matter success or failed:
//
//   x := RunSuccessAtLeast(time.Second, f)
//   x() // costs 0.1s if this attempt failed
//   x() // costs   1s if thst attempt succeeded
func RunSuccessAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Since(begin); err == nil && d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}

// RunFailedAtLeast is identical to RunAtLeast, but only failed call is ensured.
//
// Say your f() needs 0.1s no matter success or failed:
//
//   x := RunSuccessAtLeast(time.Second, f)
//   x() // costs   1s if this attempt failed
//   x() // costs 0.1s if thst attempt succeeded
func RunFailedAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Since(begin); err != nil && d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}
