package buttifier

import (
	"math/rand/v2"
	"slices"
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

type hyphenatedWord struct {
	Word        string
	Breakpoints []int
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
func (b *Buttifier) ButtifyWord(word string, breakpoints []int) (string, int) {
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

func (b *Buttifier) hyphenateWord(word string) []int {
	res := b.hyphenator.Hyphenate(word)
	if len(res) > 0 && res[len(res)-1] != len(word)-1 {
		// might not have the last syllable as a breakpoint so we add it
		res = append(res, len(word)-1)
	}

	return res
}

func (b *Buttifier) hyphenateSentence(sentence string) []*hyphenatedWord {
	words := strings.Split(sentence, " ")
	var result []*hyphenatedWord
	for _, word := range words {
		breakpoints := b.hyphenateWord(word)

		w := hyphenatedWord{
			Word:        word,
			Breakpoints: breakpoints,
		}
		result = append(result, &w)
	}
	return result
}

// replace random syllables from each word with buttWord
// returns the buttified word and true if the word was changed
func (b *Buttifier) ButtifySentence(sentence string) string {
	hyphenatedSentence := b.hyphenateSentence(sentence)
	buttifiedSyllables := 0
	totalSyllables := func() int {
		count := 0
		for _, hyphenatedWord := range hyphenatedSentence {
			count += len(hyphenatedWord.Breakpoints)
		}
		return count
	}()

	reachedButtificationRate := func() bool {
		return (float64(buttifiedSyllables) / float64(totalSyllables)) >= b.ButtificationRate
	}

	unbuttifiedWords := hyphenatedSentence
	for !reachedButtificationRate() && len(unbuttifiedWords) > 1 {
		randomWordIdx := rand.New(b.RandSource).Int() % len(unbuttifiedWords)

		buttifiedWord, buttCount := b.ButtifyWord(unbuttifiedWords[randomWordIdx].Word, unbuttifiedWords[randomWordIdx].Breakpoints)
		if buttCount > 0 {
			unbuttifiedWords[randomWordIdx].Word = buttifiedWord
			buttifiedSyllables += buttCount
			// remove the word we just buttified from the slice
			unbuttifiedWords = slices.Delete(unbuttifiedWords, randomWordIdx, randomWordIdx)
		}
	}

	result := []string{}
	for _, hyphenatedWord := range hyphenatedSentence {
		result = append(result, hyphenatedWord.Word)
	}

	return strings.Join(result, " ")
}

func (b *Buttifier) ToButtOrNotToButt() bool {
	rn := rand.New(b.RandSource).Float64()
	return rn < b.ButtificationProbability
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
