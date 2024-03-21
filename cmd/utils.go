package cmd

import (
	"fmt"
	"strings"
)

func hash_check(hash string, hash_type string) error {

	hash = strings.ToLower(hash)

	switch hash_type {

	case "SHA1":

		if len(hash) != 40 {
			return fmt.Errorf("Error: SHA1(GitHub) hash must be 40 characters long.")
		}

	case "SHA256":

		if len(hash) != 64 {
			return fmt.Errorf("Error: SHA256 hash must be 64 characters long.")
		}

	}

	for idx := 0; idx < len(hash); idx++ {

		character := hash[idx]

		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return fmt.Errorf("Error: Invalid hashsum in hash string\n%s\n%s", hash, strings.Repeat(" ", idx)+"↑")
		}

	}

	return nil

}
