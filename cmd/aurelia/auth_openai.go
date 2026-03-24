package main

import (
	"fmt"
	"io"
)

func runOpenAIAuthLogin(_ io.Reader, _ io.Writer) error {
	return fmt.Errorf("this CLI auth has been removed; use the openai provider with an API key instead")
}
