package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

func randomHex(n int) string {
	b := make([]byte, (n+1)/2)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}


func generateSlug(str string, year int) string {
	slug := strings.ToLower(str)
	slug = strings.ReplaceAll(slug, ":", "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")

	return fmt.Sprintf("%s-%d", slug, year)
}return fmt.Sprintf("%s-%d", slug, year)
}
