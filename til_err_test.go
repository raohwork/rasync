// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"testing"
)

func TestTilErr(t *testing.T) {
	myErr := errors.New("error")
	cnt := []int{0, 0}
	fgood := func() error {
		cnt[0]++
		return nil
	}
	fbad := func() error {
		cnt[1]++
		return myErr
	}

	err := TilErr(fgood, fbad, fgood)
	if err != myErr {
		t.Fatal("unexpected error: ", err)
	}
	if cnt[0] != 1 {
		t.Errorf("expect fgood ran 1 time, got %d", cnt[0])
	}
	if cnt[1] != 1 {
		t.Errorf("expect fbad ran 1 time, got %d", cnt[1])
	}
}

func TestTilErrAllDone(t *testing.T) {
	cnt := 0
	fgood := func() error {
		cnt++
		return nil
	}

	err := TilErr(fgood, fgood, fgood)
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}
	if cnt != 3 {
		t.Errorf("expect fgood ran 3 times, got %d", cnt)
	}
}

func TestTilErrAsync(t *testing.T) {
	myErr := errors.New("error")
	cnt := []int{0, 0}
	fgood := func() error {
		cnt[0]++
		return nil
	}
	fbad := func() error {
		cnt[1]++
		return myErr
	}

	err := TilErrAsync(fgood, fbad, fgood, fbad)
	if err != myErr {
		t.Fatal("unexpected error: ", err)
	}
	if cnt[0] != 2 {
		t.Errorf("expect fgood ran 2 times, got %d", cnt[0])
	}
	if cnt[1] != 2 {
		t.Errorf("expect fbad ran 2 times, got %d", cnt[1])
	}
}
