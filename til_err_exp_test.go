// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package routines

import (
	"errors"
	"fmt"
	"time"
)

func ExampleTilErr() {
	// say you have to initialize following resources in order
	//
	//   - db connection
	//   - redis connection
	//   - prefill data into redis
	//
	// instead using an init() to run them all, it's better to write
	initDB := func() error {
		fmt.Println("db")
		return nil
	}
	initRedis := func() error {
		fmt.Println("redis")
		return errors.New("cannot connect to redis")
	}
	prefill := func() error {
		fmt.Println("prefill")
		return nil
	}

	err := TilErr(initDB, initRedis, prefill)
	if err != nil {
		fmt.Println("an error occurred:", err)
	}

	// output: db
	// redis
	// an error occurred: cannot connect to redis
}

func ExampleTilErrAsync() {
	// say you have to connect to 3 independant apis
	initConn1 := func() error {
		// wait some time to simulate api call
		time.Sleep(10 * time.Millisecond)
		fmt.Println("api1")
		return nil
	}
	initConn2 := func() error {
		time.Sleep(20 * time.Millisecond)
		fmt.Println("api2")
		return errors.New("failed to connect to api2")
	}
	initConn3 := func() error {
		time.Sleep(30 * time.Millisecond)
		fmt.Println("api3")
		return errors.New("failed to connect to api3")
	}

	err := TilErrAsync(initConn1, initConn2, initConn3)
	if err != nil {
		fmt.Println("an error occurred:", err)
	}

	// output: api1
	// api2
	// api3
	// an error occurred: failed to connect to api2
}
