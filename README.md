## WIP

This project is a *work in progress*. The implementation is *incomplete* and subject to change. The documentation can be inaccurate.

# restructure

[![Build Status](https://travis-ci.org/decomp/restructure.svg?branch=master)](https://travis-ci.org/decomp/restructure)
[![Coverage Status](https://img.shields.io/coveralls/decomp/restructure.svg)](https://coveralls.io/r/decomp/restructure?branch=master)
[![GoDoc](https://godoc.org/decomp.org/x/cmd/restructure?status.svg)](https://godoc.org/decomp.org/x/cmd/restructure)

`restructure` is a tool which recovers high-level control flow primitives from control flow graphs (e.g. *.dot -> *.json). It takes an unstructured CFG (in Graphviz DOT file format) as input and produces a structured CFG (in JSON), which describes how the high-level control flow primitives relate to the nodes of the CFG.

## Installation

```shell
go get decomp.org/x/cmd/restructure
```

## Usage

```
restructure [OPTION]... CFG.dot

Flags:
  -prims string
        Comma-separated list of control flow primitives (*.dot).
  -v    Verbose output.
```

## Examples

1) Recover the high-level control flow primitives from the control flow graph [foo.dot](testdata/foo.dot).

```bash
$ restructure foo.dot
// Output:
// Isomorphism of "list" found at node "F":
//    "A"="F"
//    "B"="G"
// Isomorphism of "if" found at node "E":
//    "A"="E"
//    "B"="list0"
//    "C"="H"
```

INPUT:
* [foo.dot](testdata/foo.dot): unstructured control flow graph.

![foo.dot subgraph](https://raw.githubusercontent.com/decomp/restructure/master/testdata/foo.png)

OUTPUT:
* [foo.json](testdata/foo.json): structured control flow graph.

```c
if E {
    F
    G
}
H
```

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
