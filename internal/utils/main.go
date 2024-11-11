package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/yusuf-musleh/mmar/constants"
)

func ExtractSubdomain(host string) string {
	splitDomain := strings.Split(host, ".")
	subdomains := splitDomain[:len(splitDomain)-1]
	return strings.Join(subdomains, ".")
}

func ExtractIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)

	// Return an empty string if we could not extract IP
	if err != nil {
		return ""
	}
	return ip
}

func MmarVersionUsage() {
	fmt.Fprintf(os.Stdout, "Prints the installed version of mmar.")
}

func MmarUsage() {
	intro := `mmar is an HTTP tunnel that exposes your localhost to the world on a public URL.

Usage:
  mmar <command> [command flags]`
	fmt.Fprintln(os.Stdout, intro)

	fmt.Fprint(os.Stdout, "\nCommands:\n")

	commands := ""
	for _, subcommand := range constants.MMAR_SUBCOMMANDS {
		command := strings.Join(subcommand, "\n    ")
		commands = commands + "  " + command + "\n"
	}

	fmt.Fprintln(os.Stdout, commands)

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Run `mmar <command> -h` to get help for a specific command\n\n")
}

// Decode hash string to bytes so it can be compared
func decodeHash(hashStr string) []byte {
	dst := make([]byte, hex.DecodedLen(len([]byte(hashStr))))
	n, err := hex.Decode(dst, []byte(hashStr))
	if err != nil {
		log.Fatalf("Could not decode hash string: %v", err)
	}
	return dst[:n]
}

// Check if provided Basic Auth credentials are valid
func ValidCredentials(username string, password string) bool {
	// Compute Hash for provided username and password
	hashedUsername := sha256.Sum256([]byte(username))
	hashedPassword := sha256.Sum256([]byte(password))

	// Retrieve actual Hashes for username and password
	envHashedUsername := os.Getenv("USERNAME_HASH")
	envHashedPassword := os.Getenv("PASSWORD_HASH")

	// Compare them to check if they match and are valid
	validUsername := subtle.ConstantTimeCompare(hashedUsername[:], decodeHash(envHashedUsername)) == 1
	validPassword := subtle.ConstantTimeCompare(hashedPassword[:], decodeHash(envHashedPassword)) == 1
	return validUsername && validPassword
}
