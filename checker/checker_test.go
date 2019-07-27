package checker

import (
	"testing"
)

func TestCheckID(t *testing.T) {
	taken, err := checkID("fr3fou")

	if err != nil {
		t.Errorf("CheckIDs(\"fr3fou\") returned err = %v; expected nil", err.Error())
	} else if !taken.IsTaken {
		t.Error("CheckIDs(\"fr3fou\") returned IsTaken = false; expected true")
	}

	free, err := checkID("asdfasdfasdfasdfasdf0a9sd8f0asd8f90as8d09fa8sd09fa8s0df")

	if err != nil {
		t.Errorf("CheckIDs(\"asdfasdfasdfasdfasdf0a9sd8f0asd8f90as8d09fa8sd09fa8s0df\") returned err = %v; expected nil", err.Error())
	} else if free.IsTaken {
		t.Error("CheckIDs(\"asdfasdfasdfasdfasdf0a9sd8f0asd8f90as8d09fa8sd09fa8s0df\") returned IsTaken = true; expected false")
	}
}
