package main

import "testing"

func TestReverse(t *testing.T) {

	got := Reverse("JeffNeff")
	want := "ffeNffeJ"
	if got != want {
		t.Errorf("we got %q but wanted %q :(", got, want)
	}

}
