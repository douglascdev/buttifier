package buttifier_test

import (
	"buttifier"
	"testing"
)

type UnitTestRandSource struct{}

func (UnitTestRandSource) Uint64() uint64 {
	return 0
}

func TestButtify(t *testing.T) {
	b, err := buttifier.New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	expected, actual := "buttbuttbutt", b.Buttify("successful")
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
