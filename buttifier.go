package buttifier

import (
	"bytes"
	"math/rand/v2"
	"os"

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

// replace random syllables with buttWord and return the result
func (b *Buttifier) Buttify(word string) string {
	breakpoints := b.hyphenator.Hyphenate(word)
	// Hyphenate doesn't return the last syllable as a breakpoint so we add it
	breakpoints = append(breakpoints, len(word))
	var wordBuffer bytes.Buffer
	prev := 0
	for _, breakPoint := range breakpoints {
		rn := rand.New(b.RandSource).Float64()
		if rn < b.buttificationRate {
			wordBuffer.WriteString(b.buttWord)
		} else {
			wordBuffer.WriteString(word[prev:breakPoint])
		}
		prev = breakPoint
	}

	return wordBuffer.String()
}
