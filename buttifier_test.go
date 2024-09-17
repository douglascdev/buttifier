package buttifier

import (
	"testing"
)

type UnitTestRandSource struct{}

func (UnitTestRandSource) Uint64() uint64 {
	return 0
}

func TestButtifyWord(t *testing.T) {
	b, err := New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"successful": "buttbuttbutt",
		"someone":    "buttbutt",
	}
	for word, expected := range resultMap {
		actual, _ := b.ButtifyWord(word, b.hyphenateWord(word))
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestButtifySentence(t *testing.T) {
	b, err := New()
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"grinding for partner": "buttbutt for partner",
	}
	for sentence, expected := range resultMap {
		actual := b.ButtifySentence(sentence)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestButtifyWordKeepsCase(t *testing.T) {
	b, err := New()
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
		actual, _ := b.ButtifyWord(word, b.hyphenateWord(word))
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}
