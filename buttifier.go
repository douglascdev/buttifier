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

type syllable struct {
	Letters  string
	IdxStart int
	IdxEnd   int
}

type hyphenatedWord struct {
	Word        string
	Breakpoints []int
	Syllables   []*syllable
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
	if word == "" {
		return "", 0
	}

	var wordBuffer strings.Builder
	buttCount := 0

	for _, hyphenatedSyllable := range b.HyphenateWord(word).Syllables {
		// random float between 0 and 1
		rn := rand.New(b.RandSource).Float64()
		if rn < b.ButtificationRate {
			// normalize buttWord's case to match currentSyllable's case
			buttifiedSyllable := normalizeCase(hyphenatedSyllable.Letters, b.ButtWord)
			wordBuffer.WriteString(buttifiedSyllable)
			buttCount++
		} else {
			wordBuffer.WriteString(hyphenatedSyllable.Letters)
		}
	}

	return wordBuffer.String(), buttCount
}

func (b *Buttifier) HyphenateWord(word string) *hyphenatedWord {
	breakpoints := b.hyphenator.Hyphenate(word)

	if len(breakpoints) == 0 {
		// some words like "partne" return an empty slice, so we need to add a breakpoint
		breakpoints = []int{len(word)}
	} else if len(breakpoints) == 1 && breakpoints[0] == len(word)-1 {
		// words like "asd" return []int{2}, resulting in "as" instead of "asd"
		breakpoints[0] += 1
	} else if breakpoints[len(breakpoints)-1] != len(word) {
		// words with a single breakpoint like "partner" return []int{4}, resulting in "part" instead of "partner"
		breakpoints = append(breakpoints, len(word))
	}

	var syllables []*syllable
	idxStart := 0
	for _, breakpoint := range breakpoints {
		syllables = append(syllables, &syllable{
			Letters:  word[idxStart:breakpoint],
			IdxStart: idxStart,
			IdxEnd:   breakpoint,
		})
		idxStart = breakpoint
	}

	return &hyphenatedWord{
		Word:        word,
		Breakpoints: breakpoints,
		Syllables:   syllables,
	}
}

func (b *Buttifier) HyphenateSentence(sentence string) []*hyphenatedWord {
	words := strings.Split(sentence, " ")
	var result []*hyphenatedWord
	for _, word := range words {
		result = append(result, b.HyphenateWord(word))
	}
	return result
}

// replace random syllables from each word with buttWord
// returns the buttified word and true if the word was changed
func (b *Buttifier) ButtifySentence(sentence string) string {
	hyphenatedSentence := b.HyphenateSentence(sentence)
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

		buttifiedWord, buttCount := b.ButtifyWord(unbuttifiedWords[randomWordIdx].Word)
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
