# SPEC.md

## Overview

This system indexes a codebase into a **multi-layer dependency graph** to support:

- Agent-based code modification
- Dependency-aware planning
- Minimal context retrieval
- Safe parallel execution

The system MUST treat the codebase as:
> A graph for truth, with DAG projections for planning.

---

## Core Goals

1. Extract structural relationships from code (not just text)
2. Represent code as a queryable graph
3. Enable subgraph retrieval for tasks
4. Support DAG-based planning via SCC compression
5. Minimize token/context usage for agents

---

## Non-Goals (v1)

- Full semantic correctness across all languages
- Perfect type resolution
- Runtime tracing
- IDE-level accuracy

---

## Data Model

### Node Types

| Type | Description |
|------|-------------|
| repository | Root project |
| module | Package / namespace |
| file | Source file |
| symbol | Function, type, class, method, const |
| test | Test file or test symbol |
| build_target | Binary, lib, or build artifact |
| service | API/service boundary (optional v1) |
| config | Env/config key |

---

### Edge Types

| Edge | Description |
|------|-------------|
| contains | Hierarchy (repo → module → file → symbol) |
| imports | File/module imports another |
| exports | Symbol exported from module |
| calls | Function calls function |
| references | Symbol references symbol |
| implements | Interface/trait implementation |
| inherits | Class inheritance |
| constructs | Instantiation |
| reads | Reads field/value |
| writes | Writes field/value |
| tested_by | Symbol/file tested by test |
| builds | Build dependency |
| serves | Service dependency |
| configures | Config dependency |

---

## Graph Requirements

- Directed graph
- Stored in memory (v1) or lightweight DB (SQLite/Postgres)
- Must support:
  - Neighbor queries (1–2 hops)
  - Reverse edges (who depends on X)
  - Subgraph extraction

---

## Extraction Pipeline

### Phase 1: File Discovery

- Recursively scan repo
- Identify language per file
- Ignore:
  - `.git`
  - `node_modules`
  - build artifacts

---

### Phase 2: Parsing

Per file:

Extract:
- imports
- exports
- top-level symbols
- symbol spans (line ranges)

Optional (if feasible):
- function calls
- type references

---

### Phase 3: Graph Construction

Build:

- File → imports → file edges
- File → contains → symbols
- Symbol → calls → symbol (if resolvable)
- Symbol → references → symbol

---

### Phase 4: Indexing

Store:

Per file:
- path
- language
- hash
- imports
- symbol list

Per symbol:
- name
- kind
- signature (if available)
- file reference
- outbound edges

---

## SCC (Strongly Connected Components)

### Requirement

- Compute SCCs over symbol graph or module graph

### Purpose

- Collapse cycles into single units
- Enable DAG-based planning

### Output

- Each SCC becomes a "super-node"
- Graph becomes a DAG

---

## Planning Model

### Input

- User task / instruction

### Output

- Target symbols
- Dependency subgraph
- SCC groups
- Ordered execution plan
- Parallelizable groups
- Validation steps

---

## Subgraph Retrieval

Given a seed node:

Return:
- Node itself
- 1–2 hop neighbors
- Parent file/module
- Direct dependents
- Related tests

---

## Retrieval Strategy

### Priority Order

1. Exact match (symbol/file name)
2. Graph expansion
3. Semantic fallback (optional)

---

## Planning Rules

1. Collapse SCCs
2. Topologically sort graph
3. Identify:
   - leaf nodes (safe start)
   - high fan-in nodes (risky)
4. Partition into:
   - sequential work
   - parallel work

---

## Validation Strategy

Each plan MUST include:

- affected tests
- build targets
- runtime checks (if known)

---

## Performance Constraints

- Must handle mid-size repos (10k–100k files)
- Graph queries < 50ms
- Extraction can be slower (offline step)

---

## Storage (v1 Recommendation)

Start with:

- In-memory graph (Go structs or similar)
- Optional persistence:
  - SQLite (simple)
  - JSON snapshot

---

## Extensibility

Future additions:

- Language-specific deep parsing
- LSP integration
- Runtime traces
- Service graph
- Config graph

---

## Success Criteria

System is successful if it can:

1. Identify minimal files for a change
2. Produce correct dependency-aware plan
3. Enable parallel execution safely
4. Reduce context size significantly vs naive approaches
