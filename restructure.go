// restructure is a tool which recovers high-level control flow primitives from
// control flow graphs (e.g. *.dot -> *.json). It takes an unstructured CFG (in
// Graphviz DOT file format) as input and produces a structured CFG (in JSON),
// describes how the high-level control flow primitives relate to the nodes of
// the CFG.
//
// Example input
//
//    digraph foo {
//       E -> F
//       E -> H
//       F -> G
//       G -> H
//       E [label="entry"]
//       F
//       G
//       H [label="exit"]
//    }
//
// Example output
//
//    {
//       "entry": "if0",
//       "primitives": {
//          "if0": {
//             "primitive": "if",
//             "nodes": {
//                "A": "E",
//                "B": "list0",
//                "C": "H",
//             },
//          },
//          "list0": {
//             "primitive": "list",
//             "nodes": {
//                "A": "F",
//                "B": "G",
//             },
//          },
//       },
//    }
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	// flagPrimitives is a comma-separated list of file names to control flow
	// primitive descriptions, which are stored as Graphviz DOT files.
	flagPrimitives string
	// When flagQuiet is true, suppress non-error messages.
	flagQuiet bool
)

func init() {
	flag.StringVar(&flagPrimitives, "prims", "", "Comma-separated list of file names to control flow primitive descriptions (*.dot).")
	flag.BoolVar(&flagQuiet, "q", false, "Suppress non-error messages.")
	flag.Usage = usage
}

const use = `
restructure [OPTION]... CFG.dot
Recover control flow structures from control flow graphs (e.g. *.dot -> *.json).
`

func usage() {
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	dotPath := flag.Arg(0)
	err := restructure(dotPath)
	if err != nil {
		log.Fatalln(err)
	}
}

// restructure attempts to recover the control flow structure of a given control
// flow graph.
func restructure(dotPath string) error {
	return nil
}
