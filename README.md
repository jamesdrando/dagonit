# dagonit
Put some DAG on it

`dagonit` creates a multi-dimensional DAG representation of a codebase to support efficient analysis, task encapsulation, and task parallelization for software development workloads using agentic A.I.

---

## Table of Contents
- [About](#about)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Core Concepts](#core-concepts)
- [User Guide](#user-guide)
  - [Step 1: Indexing a Codebase](#step-1-indexing-a-codebase)
  - [Step 2: Exploring the Graph](#step-2-exploring-the-graph)
  - [Step 3: Generating an Execution Plan](#step-3-generating-an-execution-plan)
  - [Step 4: Visualizing the DAG](#step-4-visualizing-the-dag)
- [Validation Mapping](#validation-mapping)
- [Features](#features)
- [Supported Languages](#supported-languages)

---

## About
Traditional codebase analysis often treats code as a collection of text files. `dagonit` transforms this view into a queryable graph. By identifying symbols (functions, classes, etc.) and their relationships (imports, containment, testing), it allows both humans and AI agents to understand the impact of changes and plan modifications safely and efficiently.

## Prerequisites
- **Go 1.21+**: Required to build the tool.
- **Graphviz (optional)**: Required if you want to render `.dot` files into images.
- **Mermaid Live Editor (optional)**: Useful for visualizing the graph in your browser.

## Installation
Clone the repository and build the binary:

```bash
git clone https://github.com/user/dagonit.git
cd dagonit
go build -o dagonit ./cmd/dagonit/main.go
```

Move the `dagonit` binary to your `PATH` or run it from the local directory.

---

## Core Concepts

### 1. The Multi-Layer Graph
`dagonit` builds a graph where:
- **Nodes** represent files, functions, classes, and build targets.
- **Edges** represent relationships like `contains`, `imports`, `calls`, and `tested_by`.

### 2. Strongly Connected Components (SCC)
Cycles in code (e.g., File A imports File B, which imports File A) make planning difficult. `dagonit` uses Tarjan's algorithm to identify these cycles and collapse them into "SCC Super-Nodes."

### 3. The DAG Projection
By collapsing SCCs, the codebase graph becomes a **Directed Acyclic Graph (DAG)**. This allows for topological sorting, which tells you exactly which parts of the code are safe to modify first (the "leaves") and which depend on everything else (the "roots").

---

## User Guide

### Step 1: Indexing a Codebase
To begin, navigate to your project's root directory and run:

```bash
./dagonit index .
```

This scans your files, parses the symbols, and produces a `dagonit.json` file in the current directory. This file is the "brain" that all other commands use.

### Step 2: Exploring the Graph
Once indexed, you can query specific symbols or files to see their dependencies:

```bash
# Query a specific function
./dagonit query MyFunctionName

# Query a file to see its symbols and what it imports
./dagonit query path/to/file.go
```

The output will show **Out-Neighbors** (what this node depends on) and **In-Neighbors** (what depends on this node).

### Step 3: Generating an Execution Plan
If you need to refactor a specific set of symbols, use the `plan` command. It will calculate the minimal set of affected nodes and order them such that you always modify dependencies before the things that use them.

```bash
./dagonit plan SymbolA SymbolB
```

**Output:**
- **Steps**: Parallelizable groups of tasks.
- **Risks**: High fan-in nodes (nodes that many other things depend on) or detected cycles.

### Step 4: Visualizing the DAG
To get a bird's-eye view of your project's architecture, use the `serve` or `visualize` commands.

#### Interactive Web Visualizer
Start the interactive web visualizer to explore your graph with a force-directed layout:

```bash
./dagonit serve
```
This launches a server at `http://localhost:3737`. The visualizer features:
- **Interactive Force-Directed Layout:** Nodes arrange themselves naturally.
- **Dark Mode Aesthetic:** A modern, matte dark theme.
- **Pan & Zoom:** Use mouse wheel to zoom, click and drag background to pan.
- **Node Dragging:** Click and drag nodes to rearrange the graph.
- **Tooltips:** Hover over nodes to see details.
- **Physics Controls:** Adjust repulsion and spring strength in real-time.

#### Mermaid (Interactive/Browser)
Export the SCC-collapsed DAG for use in a [Mermaid Live Editor](https://mermaid.live/):

```bash
./dagonit visualize dag --format mermaid
```

#### Graphviz (Desktop)
Generate a full dependency graph as a PNG image:

```bash
./dagonit visualize full > graph.dot
dot -Tpng graph.dot -o graph.png
```

#### Plan Visualization
To see how a specific task propagates through your graph:

```bash
./dagonit visualize plan --seeds SymbolA --format mermaid
```

---

## Validation Mapping
`dagonit` automatically links source code to its validation logic:
- **Tests**: It identifies `_test.go`, `.spec.ts`, and `test_*.py` files.
- **`tested_by` Edges**: It uses heuristics to link symbols to their tests (e.g., `FuncA` is automatically linked to `TestFuncA`).
- **Build Targets**: It indexes `go.mod`, `package.json`, and `requirements.txt` as `build_target` nodes, showing you which configuration files govern which parts of the graph.

---

## Features
- **Topological Sorting**: Dependency-aware execution order.
- **Cycle Detection**: Identification of technical debt and architectural circularity.
- **Parallelization**: Detection of independent work streams.
- **Surgical Subgraphs**: Minimal context retrieval for AI agents.

## Supported Languages
- **Go**: Full AST-based parsing of imports and symbols.
- **TypeScript/JavaScript**: Regex-based parsing of imports, classes, and functions.
- **Python**: Regex-based parsing of imports, classes, and functions.
