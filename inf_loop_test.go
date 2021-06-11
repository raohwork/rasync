// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"context"
	"errors"
	"testing"
	"time"
)

type errFailed string

func (e errFailed) Error() string { return string(e) }

func TestInfiniteLoop(t *testing.T) {
	cnt := 0
	task := func() error {
		cnt++
		time.Sleep(9 * time.Millisecond)
		return nil
	}

	ctrl := InfiniteLoop(task)        // run first time
	time.Sleep(10 * time.Millisecond) // run second time during sleep
	ctrl.Cancel()

	if err := <-ctrl.Err; err != context.Canceled {
		t.Error(err)
	}

	if cnt != 2 {
		t.Fatalf("exepcted run 2 times, actually run %d times", cnt)
	}
}

func TestInfiniteLoopCancelAfterDone(t *testing.T) {
	cnt := 0
	task := func() error {
		if cnt >= 3 {
			return errFailed("jos done before cancel")
		}

		cnt++
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	ctrl := InfiniteLoop(task)
	time.Sleep(35 * time.Millisecond)
	ctrl.Cancel()

	if err := <-ctrl.Err; err == context.Canceled {
		t.Error("cancel should not work after job done")
	}
}

func TestAnyErr(t *testing.T) {
	theErr := errors.New("from failed")
	cntOk, cntFail := 0, 0
	ctrl := AnyErr(
		// always success
		InfiniteLoop(RunAtLeast(9*time.Millisecond, func() error {
			cntOk++
			return nil
		})),
		// always fail
		InfiniteLoop(RunAtLeast(10*time.Millisecond, func() error {
			cntFail++
			return theErr
		})),
	)
	defer ctrl.Cancel()

	if err := <-ctrl.Err; err != theErr {
		t.Fatal("unexpected error: ", err)
	}
	time.Sleep(10 * time.Millisecond)
	if cntOk != 2 {
		t.Errorf("expected ok to run 2 times, actually ran %d times", cntOk)
	}
	if cntFail != 1 {
		t.Errorf("expected fail to run 1 times, actually ran %d times", cntFail)
	}
}

func TestAllErr(t *testing.T) {
	theErr := errors.New("from failed")
	cntErr := []int{0, 0}
	ctrl := AllErr(
		// always fail
		InfiniteLoop(RunAtLeast(9*time.Millisecond, func() error {
			cntErr[0]++
			return theErr
		})),
		InfiniteLoop(RunAtLeast(10*time.Millisecond, func() error {
			cntErr[1]++
			return theErr
		})),
	)
	defer ctrl.Cancel()

	idx := 0
	for err := range ctrl.Err {
		idx++
		if err != theErr {
			t.Fatal("unexpected error: ", err)
		}
	}
	if idx != 2 {
		t.Fatalf("expected 2 errors, got %d", idx)
	}
	time.Sleep(10 * time.Millisecond)
	for idx, x := range cntErr {
		if x != 1 {
			t.Errorf(
				"expected %d routine to run 1 time, actually ran %d times",
				idx+1,
				x,
			)
		}
	}
}

func TestAllErrWithCancel(t *testing.T) {
	theErr := errors.New("from failed")
	cnts := []int{0, 0}
	ctrl := AllErr(
		// always fail
		InfiniteLoop(RunAtLeast(9*time.Millisecond, func() error {
			cnts[0]++
			return nil
		})),
		InfiniteLoop(RunAtLeast(10*time.Millisecond, func() error {
			cnts[1]++
			return theErr
		})),
	)
	defer ctrl.Cancel()

	if err := <-ctrl.Err; err != theErr {
		t.Fatalf("expected error returns from second routine, got %v", err)
	}
	if err := <-ctrl.Err; err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	for err := range ctrl.Err {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if cnts[0] != 2 {
		t.Errorf(
			"expected first routine to run 2 times, actually ran %d times",
			cnts[0],
		)
	}
	if cnts[1] != 1 {
		t.Errorf(
			"expected second routine to run 1 time, actually ran %d times",
			cnts[1],
		)
	}
}
