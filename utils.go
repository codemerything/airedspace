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

func generateSlug(str string, year string, check func(string) bool) string {

	slug := strings.ToLower(str)
	slug = strings.ReplaceAll(slug, ":", "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")

	if !check(slug) {
		return slug
	}
	return fmt.Sprintf("%s-%s", slug, year)

}
