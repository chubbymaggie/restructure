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
	"path/filepath"
	"sort"
	"strings"

	"decomp.org/x/graphs"
	"decomp.org/x/graphs/iso"
	"decomp.org/x/graphs/merge"
	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
)

var (
	// flagPrimitives is a comma-separated list of file names to control flow
	// primitive descriptions, which are stored as Graphviz DOT files.
	flagPrimitives string
	// When flagQuiet is true, suppress non-error messages.
	flagQuiet bool
)

func init() {
	flag.StringVar(&flagPrimitives, "prims", "", "Comma-separated list of control flow primitive descriptions (*.dot).")
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
// flow graph. It does so by repeatedly locating and merging structured
// subgraphs into single nodes until the entire graph is reduced into a single
// node or no structured subgraphs may be located.
func restructure(dotPath string) error {
	graph, err := dot.ParseFile(dotPath)
	if err != nil {
		return errutil.Err(err)
	}
	for len(graph.Nodes.Nodes) > 1 {
		prim, err := findPrim(graph)
		if err != nil {
			return errutil.Err(err)
		}
		fmt.Println("prim:", prim)
	}
	return nil
}

// A Primitive represents a high-level control flow primitive as mapping from
// the node names of the control flow primitive descriptor to the node names of
// the control flow graph.
type Primitive struct {
	// Node name of the primitive; e.g. "list0".
	NodeName string
	// Primitive name; e.g. "if", "pre_loop", ...
	PrimName string
	// Node mapping; e.g. {"A": 1, "B": 2, "C": 3}
	Nodes map[string]string
}

// findPrim locates a control flow primitive in the provided control flow graph
// and merges its nodes into a single node.
func findPrim(graph *dot.Graph) (*Primitive, error) {
	for _, sub := range subs {
		// Locate an isomorphism of sub in graph.
		m, ok := iso.Search(graph, sub)
		if !ok {
			// No match, try next control flow primitive.
			continue
		}
		printMapping(graph, sub, m)

		// Merge the nodes of the subgraph isomorphism into a single node.
		nodeName, err := merge.Merge(graph, m, sub)
		if err != nil {
			return nil, errutil.Err(err)
		}

		// Create a new control flow primitive.
		prim := &Primitive{
			NodeName: nodeName,
			PrimName: sub.Name,
			Nodes:    m,
		}
		return prim, nil
	}

	return nil, errutil.New("unable to locate control flow primitive")
}

// printMapping prints the mapping from sub node name to graph node name for an
// isomorphism of sub in graph.
func printMapping(graph *dot.Graph, sub *graphs.SubGraph, m map[string]string) {
	entry := m[sub.Entry()]
	var snames []string
	for sname := range m {
		snames = append(snames, sname)
	}
	sort.Strings(snames)
	fmt.Printf("Isomorphism of %q found at node %q:\n", sub.Name, entry)
	for _, sname := range snames {
		fmt.Printf("   %q=%q\n", sname, m[sname])
	}
}

var (
	// subs is an ordered list of subgraphs representing common control-flow
	// primitives such as 2-way conditionals, pre-test loops, etc.
	subs []*graphs.SubGraph
	// subNames specifies the name of each subgraph in subs, arranged in the same
	// order.
	subNames = []string{
		"list.dot", "if.dot", "if_else.dot", "pre_loop.dot", "post_loop.dot",
		"if_return.dot",
	}
)

func init() {
	var subPaths []string
	switch {
	case len(flagPrimitives) > 0:
		// Use custom primitives from the comma-separated list in the "-prims"
		// flag.
		subPaths = strings.Split(flagPrimitives, ",")
	default:
		// Use default primitives.
		subDir, err := goutil.SrcDir("decomp.org/x/graphs/testdata/primitives")
		if err != nil {
			log.Fatalln(errutil.Err(err))
		}
		for _, subName := range subNames {
			subPath := filepath.Join(subDir, subName)
			subPaths = append(subPaths, subPath)
		}
	}

	// Parse subgraphs representing control flow primitives.
	for _, subPath := range subPaths {
		sub, err := graphs.ParseSubGraph(subPath)
		if err != nil {
			log.Fatalln(errutil.Err(err))
		}
		subs = append(subs, sub)
	}
}
