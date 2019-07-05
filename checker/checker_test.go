package checker

import (
	"sync"
	"testing"
)

func TestCheckIDWithAPI(t *testing.T) {
	res := make(chan SteamID, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	checkID("fr3fou", &wg, res)
	wg.Wait()

	taken := <-res
	if !taken.IsTaken {
		t.Error("CheckIDs(\"fr3fou\", &wg, res) returned IsTaken = false; expected true")
	}

}
