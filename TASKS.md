# TASKS.md

## Phase 0: Project Setup

- [ ] Create project structure
- [ ] Define core packages/modules:
  - graph
  - parser
  - indexer
  - planner
  - retrieval

---

## Phase 1: File Discovery

- [ ] Implement recursive file scanner
- [ ] Add ignore rules:
  - .git
  - node_modules
  - dist/build dirs
- [ ] Detect language by extension

---

## Phase 2: Basic Parsing (v1 = shallow)

For each supported language:

- [ ] Extract imports
- [ ] Extract top-level symbols
- [ ] Record file metadata

Start with:
- Go
- TypeScript/JavaScript
- Python

Fallback:
- Regex-based parsing is acceptable for v1

---

## Phase 3: Graph Core

### Data Structures

- [ ] Define Node struct
- [ ] Define Edge struct
- [ ] Define Graph struct

Graph must support:
- [ ] Add node
- [ ] Add edge
- [ ] Get neighbors
- [ ] Get reverse neighbors

---

## Phase 4: Graph Construction

- [ ] Build file-level graph (imports)
- [ ] Build containment graph (file → symbols)
- [ ] Build basic symbol references (if possible)

---

## Phase 5: Index Storage

- [ ] Store file metadata
- [ ] Store symbol metadata
- [ ] Link nodes and edges

Optional:
- [ ] Serialize graph to JSON
- [ ] Load graph from JSON

---

## Phase 6: SCC Computation

- [ ] Implement Tarjan’s algorithm (or Kosaraju)
- [ ] Compute SCCs over:
  - module graph (minimum)
  - symbol graph (optional)

- [ ] Collapse SCCs into super-nodes
- [ ] Produce DAG

---

## Phase 7: Retrieval Engine

### Input:
- symbol name / file name

### Output:
- subgraph

Tasks:

- [ ] Exact match lookup
- [ ] 1-hop expansion
- [ ] 2-hop expansion (optional)
- [ ] Include:
  - parent file
  - dependents
  - related tests

---

## Phase 8: Planning Engine

Given:
- user task
- seed nodes

Produce:

- [ ] Identify affected nodes
- [ ] Expand dependency neighborhood
- [ ] Collapse SCCs
- [ ] Topologically sort DAG

Output:

- ordered steps
- parallel groups
- risk nodes

---

## Phase 9: Validation Mapping

- [ ] Link files/symbols to tests
- [ ] Link to build targets (if detectable)

---

## Phase 10: CLI Interface (v1)

Commands:

- [ ] index
- [ ] query <symbol>
- [ ] plan <task>

---

## Phase 11: Optimization

- [ ] Cache graph
- [ ] Incremental indexing (optional)
- [ ] Memory optimization

---

## Phase 12: Testing

- [ ] Small repo test
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
- [ ] Visualization (graph output)

---

## Definition of Done

System can:

- [ ] Index a repo
- [ ] Retrieve dependency subgraph
- [ ] Produce ordered plan
- [ ] Identify parallelizable work
- [ ] Suggest validation steps
