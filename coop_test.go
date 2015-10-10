// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package coop

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestAt_Future(test *testing.T) {
	start := time.Now()
	done := At(start.Add(100*time.Millisecond), func() {
		diff := time.Now().Sub(start)
		if diff > 105*time.Millisecond {
			test.Errorf("Expected to run in 100 ms, it did in %v", diff)
		}
	})
	<-done
}

func TestAt_Past(test *testing.T) {
	start := time.Now()
	done := At(start.Add(-100*time.Millisecond), func() {})
	<-done
	diff := time.Now().Sub(start)
	if diff > time.Millisecond {
		test.Errorf("Expected to return immediately, but it took %v", diff)
	}
}

func TestAfter_Future(test *testing.T) {
	start := time.Now()
	done := After(100*time.Millisecond, func() {
		diff := time.Now().Sub(start)
		if diff > 105*time.Millisecond {
			test.Errorf("Expected to run in 100 ms, it did in %v", diff)
		}
	})
	<-done
}

func TestEvery(test *testing.T) {
	dur := 10 * time.Millisecond
	count := 0
	Every(dur, func() {
		count++
	})
	<-time.After(100 * time.Millisecond)
	count++
	if count < 9 {
		test.Errorf("Expected to run in at least 9 times, it did %v times", count)
	}
}

func TestUntil_Future(test *testing.T) {
	count := 0
	done := Until(time.Now().Add(100*time.Millisecond), 10*time.Millisecond, func() {
		count++
	})
	<-done
	if count < 9 {
		test.Errorf("Expected to run for at least for 9 times, but it ran for %v times", count)
	}
}

func TestUntil_Past(test *testing.T) {
	count := 0
	done := Until(time.Now().Add(-100*time.Millisecond), 10*time.Millisecond, func() {
		count++
	})
	<-done
	if count != 0 {
		test.Errorf("Expected to run for at least for 0 times, but it ran for %v times", count)
	}
}

func TestTimeout_TimedOut(test *testing.T) {
	done := Timeout(100*time.Millisecond, func() {
		time.Sleep(time.Minute)
	})
	if <-done {
		test.Errorf("Expected to get timed out, but it has been completed")
	}
}

func TestTimeout_Completed(test *testing.T) {
	done := Timeout(time.Minute, func() {
		time.Sleep(100 * time.Millisecond)
	})
	if !<-done {
		test.Errorf("Expected to get completed, but it has been timed out")
	}
}

func TestAll(test *testing.T) {
	var ndone int32 = 0
	done := All(func() {
		atomic.AddInt32(&ndone, 1)
	}, func() {
		atomic.AddInt32(&ndone, 1)
	}, func() {
		atomic.AddInt32(&ndone, 1)
	})
	<-done
	if atomic.LoadInt32(&ndone) != 3 {
		test.Errorf("Only %d fn completed, expected 3", ndone)
	}
}

func TestAllWithThrottle(test *testing.T) {
	var nlive int32 = 0
	var ndone int32 = 0
	sch := make(chan struct{}, 5)
	cch := make(chan struct{})
	fch := make(chan struct{})
	fn := func() {
		atomic.AddInt32(&nlive, 1)
		sch <- struct{}{} // Notify that nlive has been updated
		<-cch             // Wait for permission to continue
		atomic.AddInt32(&ndone, 1)
		fch <- struct{}{}
	}
	done := AllWithThrottle(2, fn, fn, fn, fn, fn)
	<-sch
	<-sch
	// If less than 2 fn are started, the test will hang.

	l1 := atomic.LoadInt32(&nlive)
	if l1 != 2 {
		test.Errorf("Expected 2 live fn, got %d", l1)
	}

	// Let 2 fn continue
	cch <- struct{}{}
	cch <- struct{}{}

	// Wait for 2 fn to finish and update ndone
	<-fch
	<-fch

	// Wait for two more fn to start and update nlive
	<-sch
	<-sch

	l2 := atomic.LoadInt32(&nlive)
	if l2 != 4 {
		test.Errorf("Expected 4 live fn, got %d", l2)
	}

	d1 := atomic.LoadInt32(&ndone)
	if d1 != 2 {
		test.Errorf("Expected 2 done fn, got %d", d1)
	}

	// Let 2 fn continue
	cch <- struct{}{}
	cch <- struct{}{}

	// Wait for 2 fn to finish and update ndone
	<-fch
	<-fch

	// Wait for last more fn to start and update nlive
	<-sch

	l3 := atomic.LoadInt32(&nlive)
	if l3 != 5 {
		test.Errorf("Expected 5 live fn, got %d", l3)
	}

	d2 := atomic.LoadInt32(&ndone)
	if d2 != 4 {
		test.Errorf("Expected 4 done fn, got %d", d2)
	}

	// Let last fn continue
	cch <- struct{}{}

	// Wait for last fn to finish and update ndone
	<-fch

	<-done

	d3 := atomic.LoadInt32(&ndone)
	if d3 != 5 {
		test.Errorf("Expected 5 done fn, got %d", d3)
	}
}

func TestReplicate(test *testing.T) {
	results := make(chan bool, 5)
	done := Replicate(5, func() {
		results <- true
	})
	<-done
	close(results)
	count := 0
	for _ = range results {
		count++
	}
	if count != 5 {
		test.Errorf("Expected 5 to run, but %v worked", count)
	}
}
