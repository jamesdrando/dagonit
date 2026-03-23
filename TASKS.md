# TASKS.md

## Phase 0: Project Setup

- [x] Create project structure
- [x] Define core packages/modules:
  - graph
  - parser
  - indexer
  - planner
  - retrieval

---

## Phase 1: File Discovery

- [x] Implement recursive file scanner
- [x] Add ignore rules:
  - .git
  - node_modules
  - dist/build dirs
- [x] Detect language by extension

---

## Phase 2: Basic Parsing (v1 = shallow)

For each supported language:

- [x] Extract imports
- [x] Extract top-level symbols
- [x] Record file metadata

Start with:
- [x] Go
- [x] TypeScript/JavaScript
- [x] Python

---

## Phase 3: Graph Core

### Data Structures

- [x] Define Node struct
- [x] Define Edge struct
- [x] Define Graph struct

Graph must support:
- [x] Add node
- [x] Add edge
- [x] Get neighbors
- [x] Get reverse neighbors

---

## Phase 4: Graph Construction

- [x] Build file-level graph (imports)
- [x] Build containment graph (file → symbols)
- [x] Build basic symbol references (if possible)

---

## Phase 5: Index Storage

- [x] Store file metadata
- [x] Store symbol metadata
- [x] Link nodes and edges

Optional:
- [x] Serialize graph to JSON
- [x] Load graph from JSON

---

## Phase 6: SCC Computation

- [x] Implement Tarjan’s algorithm (or Kosaraju)
- [x] Compute SCCs over:
  - module graph (minimum)
  - symbol graph (optional)

- [x] Collapse SCCs into super-nodes
- [x] Produce DAG

---

## Phase 7: Retrieval Engine

### Input:
- symbol name / file name

### Output:
- subgraph

Tasks:

- [x] Exact match lookup
- [x] 1-hop expansion
- [x] 2-hop expansion (optional)
- [x] Include:
  - parent file
  - dependents
  - related tests

---

## Phase 8: Planning Engine

Given:
- user task
- seed nodes

Produce:

- [x] Identify affected nodes
- [x] Expand dependency neighborhood
- [x] Collapse SCCs
- [x] Topologically sort DAG

Output:

- [x] ordered steps
- [x] parallel groups
- [x] risk nodes

---

## Phase 9: Validation Mapping

- [x] Link files/symbols to tests
- [x] Link to build targets (if detectable)
- [x] Implement tested_by edge heuristic

---

## Phase 10: CLI Interface (v1)

Commands:

- [x] index
- [x] query <symbol>
- [x] plan <task>

---

## Phase 11: Visualization

- [x] Export graph to DOT (Graphviz)
- [x] Export graph to Mermaid (MD/HTML)
- [x] Visualize SCC-collapsed DAG
- [x] Highlight planning steps in graph

---

## Phase 12: Optimization

- [ ] Cache graph
- [ ] Incremental indexing (optional)
- [ ] Memory optimization

---

## Phase 13: Testing

- [x] Small repo test
- [ ] Medium repo test
- [ ] Validate:
  - graph correctness
  - SCC correctness
  - planning output sanity

---

## Stretch Goals

- [ ] LSP integration
- [ ] Better call graph resolution
- [ ] Language plugins

---

## Definition of Done

System can:

- [x] Index a repo
- [x] Retrieve dependency subgraph
- [x] Produce ordered plan
- [x] Identify parallelizable work
- [x] Suggest validation steps
- [x] Visualize the DAG and execution plan
