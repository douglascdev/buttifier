package buttifier

import (
	"math/rand/v2"
	"os"
	"strings"

	"github.com/speedata/hyphenation"
)

/*
implements rand.Source
replaced in unit tests to pass a custom random seed
*/
type DefaultRandSource struct{}

func (DefaultRandSource) Uint64() uint64 {
	return rand.Uint64()
}

type Buttifier struct {
	buttWord                 string
	hyphenator               *hyphenation.Lang
	buttificationProbability float64
	buttificationRate        float64
	RandSource               rand.Source
}

func New() (*Buttifier, error) {
	hyph, err := newHyphenator("hyph-en-us.pat.txt")
	if err != nil {
		return nil, err
	}
	return &Buttifier{
		buttWord:                 "butt",
		hyphenator:               hyph,
		buttificationProbability: 0.1,
		buttificationRate:        0.3,
		RandSource:               DefaultRandSource{},
	}, nil
}

func NewWithCustomButt(
	buttWord string,
	hyphenationFile string,
	buttificationProbability float64,
	buttificationRate float64,
) (*Buttifier, error) {
	hyph, err := newHyphenator("hyph-en-us.pat.txt")
	if err != nil {
		return nil, err
	}
	return &Buttifier{
		buttWord:                 buttWord,
		hyphenator:               hyph,
		buttificationProbability: buttificationProbability,
		buttificationRate:        buttificationRate,
		RandSource:               DefaultRandSource{},
	}, nil
}

func newHyphenator(hyphenationFile string) (*hyphenation.Lang, error) {
	r, err := os.Open(hyphenationFile)
	if err != nil {
		return nil, err
	}
	return hyphenation.New(r)
}

// replace random syllables with buttWord
// returns the buttified word and the number of buttified syllables
func (b *Buttifier) ButtifyWord(word string) (string, int) {
	breakpoints := b.hyphenator.Hyphenate(word)
	// Hyphenate doesn't return the last syllable as a breakpoint so we add it
	breakpoints = append(breakpoints, len(word))
	var wordBuffer strings.Builder
	prev := 0
	buttCount := 0
	for _, breakPoint := range breakpoints {
		rn := rand.New(b.RandSource).Float64()
		if rn < b.buttificationRate {
			wordBuffer.WriteString(b.buttWord)
			buttCount++
		} else {
			wordBuffer.WriteString(word[prev:breakPoint])
		}
		prev = breakPoint
	}

	return wordBuffer.String(), buttCount
}

// replace random syllables from each word with buttWord
// returns the buttified word and true if the word was changed
func (b *Buttifier) ButtifySentence(sentence string) (string, bool) {
	rn := rand.New(b.RandSource).Float64()
	toButtOrNotToButt := rn < b.buttificationProbability

	if !toButtOrNotToButt {
		return sentence, false
	}

	words := strings.Split(sentence, " ")
	changed := false
	for i := range words {
		buttifiedWord, count := b.ButtifyWord(words[i])
		if count > 0 {
			words[i] = buttifiedWord
			changed = true
		}
	}
	return strings.Join(words, " "), changed
}
