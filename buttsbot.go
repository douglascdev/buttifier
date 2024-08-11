package buttsbot

import (
	"math/rand"
	"os"

	"github.com/speedata/hyphenation"
)

type Buttifier struct {
	buttWord                 string
	hyphenator               *hyphenation.Lang
	buttificationProbability float64
	buttificationRate        float64
	rand                     *rand.Rand
}

func New() (*Buttifier, error) {
	hyph, err := newHyphenator("hyph-en-us.pat.txt")
	if err != nil {
		return nil, err
	}
	return &Buttifier{buttWord: "butt", hyphenator: hyph, buttificationProbability: 0.1, buttificationRate: 0.3, rand: rand.New(rand.SourNewSource(seed))}, nil
}

func NewWithCustomButt(buttWord string, hyphenationFile string, buttificationProbability float64, buttificationRate float64) (*Buttifier, error) {
	hyph, err := newHyphenator("hyph-en-us.pat.txt")
	if err != nil {
		return nil, err
	}
	return &Buttifier{buttWord: buttWord, hyphenator: hyph}, nil
}

func newHyphenator(hyphenationFile string) (*hyphenation.Lang, error) {
	r, err := os.Open(hyphenationFile)
	if err != nil {
		return nil, err
	}
	return hyphenation.New(r)
}

func (b *Buttifier) Buttify(word string) string {
	breakpoints := b.hyphenator.Hyphenate(word)
	for i, breakPoint := range breakpoints {
		if rand.Float64() < b.buttificationRate {
			start := 0
			if i > 0 {
				start = breakpoints[i-1]
			} else {
				start = 0
			}
			word = word[0:start] + b.buttWord + word[breakPoint:]
		}
	}

	return word
}
