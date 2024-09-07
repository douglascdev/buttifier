package buttifier

import (
	"math/rand/v2"
	"strings"
	"unicode"

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
	hyphenator               *hyphenation.Lang
	ButtWord                 string
	ButtificationProbability float64
	ButtificationRate        float64
	RandSource               rand.Source
}

func New() (*Buttifier, error) {
	hyph, err := hyphenation.New(strings.NewReader(HyphenatorData))
	if err != nil {
		return nil, err
	}
	return &Buttifier{
		ButtWord:                 "butt",
		hyphenator:               hyph,
		ButtificationProbability: 0.05,
		ButtificationRate:        0.3,
		RandSource:               DefaultRandSource{},
	}, nil
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
		// random float between 0 and 1
		rn := rand.New(b.RandSource).Float64()
		currentSyllable := word[prev:breakPoint]
		if rn < b.ButtificationRate {
			// normalize buttWord's case to match currentSyllable's case
			buttifiedSyllable := normalizeCase(currentSyllable, b.ButtWord)
			wordBuffer.WriteString(buttifiedSyllable)
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
	toButtOrNotToButt := rn < b.ButtificationProbability

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

// tries to normalize buttWord's case to match currentSyllable's case
// ("SOMeone", "buttbutt") -> "BUTtbutt"
// ("SOMEone", "buttbutt") -> "BUTTbutt"
func normalizeCase(currentSyllable string, buttWord string) string {
	buttifiedSyllable := strings.Split(buttWord, "")

	if len(currentSyllable) < len(buttWord) {
		isAllUpperCase := true
		for i := 0; i < len(currentSyllable); i++ {
			if !unicode.IsUpper(rune(currentSyllable[i])) {
				isAllUpperCase = false
				break
			}
		}
		if isAllUpperCase {
			return strings.ToUpper(strings.Join(buttifiedSyllable, ""))
		}
	}

	// copies case character by character
	for i := 0; i < min(len(buttWord), len(currentSyllable)); i++ {
		letterIsUpperCase := unicode.IsUpper(rune(currentSyllable[i]))
		if letterIsUpperCase {
			buttifiedSyllable[i] = strings.ToUpper(buttifiedSyllable[i])
		} else {
			buttifiedSyllable[i] = strings.ToLower(buttifiedSyllable[i])
		}
	}

	return strings.Join(buttifiedSyllable, "")
}
