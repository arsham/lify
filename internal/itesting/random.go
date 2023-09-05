// Package itesting provides some helper functions for testing.
package itesting

import "math/rand"

// RandomDNA returns a randomly generates string with the length of count.
func RandomDNA(count int) string {
	const runes = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, count)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
