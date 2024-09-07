# buttifier

<p align="center">
  <img src="https://cdn.7tv.app/emote/664d082bba453b48ddc926a2/4x.gif" alt="butt" />
</p>

Replace syllables with `butt`(default) or a custom word. Made for a twitch bot.

## Installation

```bash
go get github.com/douglascdev/buttifier
```

## Usage

```go
package main

import "github.com/douglascdev/buttifier"

func main() {
	buttifier, err := buttifier.New()

	// 50% chance of buttifying the sentence passed to ButtifySentence
	buttifier.ButtificationProbability = 0.5
	// buttify about 30% of the syllables
	buttifier.ButtificationRate = 0.3
	// what each buttified syllable should be replaced with
	buttifier.ButtWord = "butt"

	if err != nil {
		panic(err)
	}
	newSentence, didButtify := buttifier.ButtifySentence("Someone did that something something")
	if didButtify {
		// Someone butt that something something
		println(newSentence)
		return
	}
	println("Did not buttify sentence")
}
```
