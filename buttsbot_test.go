package buttsbot_test

import (
	"buttsbot"
	"math/rand"
	"testing"
)

func TestButtify(t *testing.T) {
	rand.New(rand.NewSource(42))

	b, err := buttsbot.New()
	if err != nil {
		t.Fatal(err)
	}
	if b.Buttify("Contributor") != "Conbuttutor" {
		t.Errorf("Expected %s, got %s", "Conbuttutor", b.Buttify("Contributor"))
	}
}
