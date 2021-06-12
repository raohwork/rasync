// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

// TilErr iterates funcs until first error
//
// It guarantees f in funcs is executed one by one in order, and stops at first
// error. See example for common use case.
func TilErr(funcs ...func() error) (err error) {
	for _, f := range funcs {
		if err = f(); err != nil {
			return
		}
	}

	return
}

// TilErrAsync runs all funcs in background, block til done and return first error
//
// It guarantees all f in funcs are executed. Only first error is returned, others
// are ignored. See example for common use case.
func TilErrAsync(funcs ...func() error) (err error) {
	ch := make(chan error, len(funcs))
	defer close(ch)

	for _, f := range funcs {
		go func(f func() error) {
			ch <- f()
		}(f)
	}

	for x := 0; x < len(funcs); x++ {
		e := <-ch
		if err != nil || e == nil {
			continue
		}
		err = e
	}

	return
}
