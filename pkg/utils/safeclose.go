package utils

import (
	"io"
	"log"
)

func SafeClose(c io.Closer) {
	if c != nil {
		if err := c.Close(); err != nil {
			log.Printf("error closing: %v", err)
		}
	}
}
