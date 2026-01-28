package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	passwords := map[string]string{
		"admin123":   "admin",
		"petugas123": "petugas",
	}

	for pass, user := range passwords {
		hash := sha256.Sum256([]byte(pass))
		fmt.Printf("User: %s | Password: %s | Hash: %x\n", user, pass, hash)
	}
}
