package main

import (
	"flag"
	"fmt"
	"giftlock/internal/security"
	"os"
)

func main() {
	password := flag.String("password", "", "Password to hash")
	help := flag.Bool("help", false, "Show usage information")
	flag.Parse()
	if *help {
		fmt.Println("Usage: generate_password --password <yourpassword>")
		fmt.Println("Example: go run scripts/generate_password.go --password \"hunter2\"")
		fmt.Println("Returns the hashed password.")
		os.Exit(0)
	}
	if *password == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --password argument is required")
		fmt.Fprintln(os.Stderr, "Use --help for usage information.")
		os.Exit(1)
	}
	h, err := security.HashPassword(*password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
	fmt.Println(h)
}
