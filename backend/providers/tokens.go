package providers

import (
	"github.com/tiktoken-go/tokenizer"
)

// TokenCount calculates the number of tokens in the given text using the GPT-4o tokenizer.
func TokenCount(text string) (int, error) {
	enc, err := tokenizer.Get(tokenizer.O200kBase)
	if err != nil {
		return 0, err
	}

	return enc.Count(text)
}
