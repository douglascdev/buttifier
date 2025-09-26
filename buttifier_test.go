package buttifier

import (
	"fmt"
	"strings"
	"testing"
)

type UnitTestRandSource struct{}

func (UnitTestRandSource) Uint64() uint64 {
	return 0
}

func TestButtifyWord(t *testing.T) {
	b, err := New(English)
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"successful": "buttbuttbutt",
		"someone":    "buttbutt",
		"":           "",
		"partner":    "buttbutt",
		"partne":     "butt",
		"developers": "buttbuttbuttbuttbutt",
		"computer":   "buttbuttbutt",
		"asd":        "butt",
		"butt":       "butt",
	}
	for word, expected := range resultMap {
		actual, _ := b.ButtifyWord(word)
		if expected != actual {
			t.Errorf("expected %s => %s, got %s => %s", word, expected, word, actual)
		}
	}
}

func TestButtifySentence(t *testing.T) {
	b, err := New(English)
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}
	resultMap := map[string]string{
		"grinding for partner":     "buttbutt for partner",
		"frizze5Wade laffer curve": "butt laffer curve",
	}
	for sentence, expected := range resultMap {
		actual := b.ButtifySentence(sentence)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestButtifyWordKeepsCase(t *testing.T) {
	b, err := New(English)
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
		"asd":        "butt",
	}
	for word, expected := range resultMap {
		actual, _ := b.ButtifyWord(word)
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestHyphenateWord(t *testing.T) {
	b, err := New(English)
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}

	expectedResults := []string{
		"successful",
		"someone",
		"",
		"partner",
		"partne",
		"developers",
		"computer",
		"asd",
		"butt",
	}
	for _, expected := range expectedResults {
		hyphenatedWord := b.HyphenateWord(expected)
		syllables := []string{}
		for _, syllable := range hyphenatedWord.Syllables {
			syllables = append(syllables, syllable.Letters)
		}

		actual := strings.Join(syllables, "")
		if expected != actual {
			t.Errorf(fmt.Sprintf("expected '%s' got '%s'", expected, actual))
		}
	}
}

func TestHyphenateWordPt(t *testing.T) {
	b, err := New(Portuguese)
	b.RandSource = UnitTestRandSource{}
	if err != nil {
		t.Fatal(err)
	}

	expectedResults := map[string][]string{
		"perspicaz":  {"pers", "pi", "caz"},
		"milionario": {"mi", "li", "o", "na", "ri", "o"},
	}
	for word, expectedSyllables := range expectedResults {
		hyphenatedWord := b.HyphenateWord(word)

		syllables := []string{}
		for _, syllable := range hyphenatedWord.Syllables {
			syllables = append(syllables, syllable.Letters)
		}

		if len(expectedSyllables) != len(hyphenatedWord.Syllables) {
			t.Errorf("expected '%d' syllables got '%d' for word %q", len(expectedSyllables), len(hyphenatedWord.Syllables), word)
		}

		for i, syllable := range syllables {
			if syllable != expectedSyllables[i] {
				t.Errorf(fmt.Sprintf("expected '%s' got '%s'", strings.Join(expectedSyllables, "-"), strings.Join(syllables, "-")))
			}
		}

	}
}
