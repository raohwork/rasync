// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"testing"
)

func TestRecorded(t *testing.T) {
	idx := uint64(0)
	f := func(i uint64) error {
		idx = i
		return nil
	}
	g := Recorded(f)

	for x := uint64(0); x < 10; x++ {
		g()
		if idx != x {
			t.Fatalf("expected %d, got %d", x, idx)
		}
	}
}

func TestRetry(t *testing.T) {
	theErr := errors.New("the error")
	f := func(i uint64) error {
		if i >= 5 {
			return nil
		}
		return theErr
	}
	ch := Retry(Recorded(f))
	cnt := 0
	for err := range ch {
		cnt++
		if err != theErr {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if cnt != 5 {
		t.Fatalf("expected to failed 5 times, got %d", cnt)
	}
}

func TestTriesAtMostFailed(t *testing.T) {
	theErr := errors.New("the error")
	f := func(i uint64) error {
		if i >= 5 {
			return nil
		}
		return theErr
	}
	ch := TriesAtMost(3, Recorded(f))
	cnt := 0
	for err := range ch {
		cnt++
		if err != theErr {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if cnt != 3 {
		t.Fatalf("expected to failed 3 times, got %d", cnt)
	}
}

func TestTriesAtMostDone(t *testing.T) {
	theErr := errors.New("the error")
	f := func(i uint64) error {
		if i >= 3 {
			return nil
		}
		return theErr
	}
	ch := TriesAtMost(5, Recorded(f))
	cnt := 0
	for err := range ch {
		cnt++
		if err != theErr {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if cnt != 3 {
		t.Fatalf("expected to failed 3 times, got %d", cnt)
	}
}

func TestTryAtMostFailed(t *testing.T) {
	theErr := errors.New("the error")
	f := func(i uint64) error {
		if i >= 3 {
			return nil
		}
		return theErr
	}
	err := TryAtMost(3, Recorded(f))
	if err == nil {
		t.Fatal("expected error, got nothing")
	}
	if err != theErr {
		t.Fatal("unexpected error: ", err)
	}
}

func TestTryAtMostDone(t *testing.T) {
	theErr := errors.New("the error")
	f := func(i uint64) error {
		if i >= 3 {
			return nil
		}
		return theErr
	}
	err := TryAtMost(5, Recorded(f))
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}
}
