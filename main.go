package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI with args (excluding the program name) and returns the exit code.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 || args[0] != "get" {
		fmt.Fprintln(stderr, "usage: go-cli-pokedex get <pokemon> [pokemon...]")
		return 1
	}

	code := 0
	emitted := 0 // results or errors shown so far; each new one gets a blank line before it
	for _, arg := range args[1:] {
		name := strings.ToLower(arg)
		p, err := FetchPokemon(name)
		if err != nil {
			if emitted > 0 {
				fmt.Fprintln(stderr)
			}
			emitted++
			fmt.Fprintf(stderr, "error: %s: %v\n", name, err)
			code = 1
			continue
		}
		if emitted > 0 {
			fmt.Fprintln(stdout)
		}
		emitted++
		fmt.Fprintf(stdout, "Name: %s\nWeight: %.1f kg\nTypes: %s\n", p.Name, p.WeightKg, strings.Join(p.Types, ", "))
	}
	return code
}
