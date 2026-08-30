// Command lo is the real Phase 1 LO compiler CLI: lex -> parse -> emit real `.prn` text.
// Usage: lo build <in.llll> -o <out.prn>
package main

import (
	"fmt"
	"os"

	"github.com/emilyspringerton/LO/internal/emitter"
	"github.com/emilyspringerton/LO/internal/lexer"
	"github.com/emilyspringerton/LO/internal/parser"
)

func main() {
	if len(os.Args) < 5 || os.Args[1] != "build" || os.Args[3] != "-o" {
		fmt.Fprintln(os.Stderr, "usage: lo build <in.llll> -o <out.prn>")
		os.Exit(1)
	}
	inPath, outPath := os.Args[2], os.Args[4]

	src, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lo: %v\n", err)
		os.Exit(1)
	}

	toks, err := lexer.Lex(string(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "lo: %v\n", err)
		os.Exit(1)
	}
	prog, err := parser.Parse(toks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lo: %v\n", err)
		os.Exit(1)
	}
	out, err := emitter.Emit(prog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lo: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "lo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("lo: %s -> %s\n", inPath, outPath)
}
