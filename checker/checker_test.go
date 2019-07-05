package checker

import (
	"sync"
	"testing"
)

func TestCheckID(t *testing.T) {
	res := make(chan SteamID, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	checkID("fr3fou", &wg, res)

	taken := <-res
	if !taken.IsTaken {
		t.Error("CheckIDs(\"fr3fou\", &wg, res) returned IsTaken = false; expected true")
	}

	wg.Add(1)
	checkID("asdfasdfasdfasdfasdf0a9sd8f0asd8f90as8d09fa8sd09fa8s0df", &wg, res)

	taken = <-res
	if taken.IsTaken {
		t.Error("CheckIDs(\"asdfasdfasdfasdfasdf0a9sd8f0asd8f90as8d09fa8sd09fa8s0df\", &wg, res) returned IsTaken = true; expected false")
	}
}
