package buttifier_test

import (
	"testing"

	"github.com/douglascdev/buttifier"
)

type UnitTestRandSource struct{}

func (UnitTestRandSource) Uint64() uint64 {
	return 0
}

func TestButtifyWord(t *testing.T) {
	b, err := buttifier.New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"successful": "buttbuttbutt",
		"someone":    "buttbutt",
	}
	for word, expected := range resultMap {
		actual, _ := b.ButtifyWord(word)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestButtifySentence(t *testing.T) {
	b, err := buttifier.New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"grinding for partner": "buttbutt butt buttbutt",
	}
	for sentence, expected := range resultMap {
		actual, _ := b.ButtifySentence(sentence)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestButtifyWordKeepsCase(t *testing.T) {
	b, err := buttifier.New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"SUCCESSFUL": "BUTTBUTTBUTT",
		"Someone":    "Buttbutt",
		"SOmeone":    "BUttbutt",
		"SOMeone":    "BUTtbutt",
		"SOMEone":    "BUTTbutt",
	}
	for word, expected := range resultMap {
		actual, _ := b.ButtifyWord(word)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}
