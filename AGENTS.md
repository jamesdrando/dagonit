# AGENTS.md

## Purpose

This system is designed to support AI agents performing code modifications safely and efficiently.

Agents MUST operate using the graph, not raw file guessing.

---

## Core Principles

1. **Graph is truth**
2. **Tree is presentation**
3. **Always work on subgraphs**
4. **Never assume independence without checking dependencies**
5. **Collapse cycles before planning**

---

## Required Workflow

### Step 1: Understand Task

Agent MUST:

- Restate goal in terms of code changes
- Identify likely symbols/files

---

### Step 2: Locate Entry Points

Use:

- exact symbol match
- file name match

DO NOT use semantic search first.

---

### Step 3: Expand Subgraph

Retrieve:

- node
- neighbors (1–2 hops)
- parent file/module
- dependents
- related tests

---

### Step 4: Analyze Dependencies

Agent MUST determine:

- upstream dependencies (what this uses)
- downstream dependencies (what uses this)

---

### Step 5: Detect Cycles

- Identify SCCs
- Treat each SCC as a single unit

NEVER attempt to partially modify a cycle without understanding full scope.

---

### Step 6: Plan Execution

Agent MUST:

1. Collapse SCCs
2. Topologically sort
3. Identify:
   - safe starting points (leaves)
   - risky nodes (high fan-in)

---

### Step 7: Identify Parallel Work

Parallelizable if:

- nodes do not depend on each other
- not in same SCC

---

### Step 8: Generate Plan

Plan MUST include:

- ordered steps
- parallel groups
- affected files
- affected symbols
- validation steps

---

### Step 9: Validate Before Modify

Agent MUST:

- identify tests
- identify build targets
- confirm assumptions

---

### Step 10: Execute Changes

Rules:

- make minimal changes
- preserve interfaces when possible
- introduce compatibility layers if needed

---

## Anti-Patterns (FORBIDDEN)

- Editing files without dependency analysis
- Using only semantic search to find code
- Ignoring downstream impact
- Modifying shared utilities blindly
- Breaking cycles incorrectly

---

## Heuristics

### Safe Nodes

- leaf nodes
- low fan-in
- internal/private symbols

### Risky Nodes

- high fan-in
- shared utilities
- core types/interfaces

---

## Retrieval Rules

Priority:

1. Exact match
2. Graph expansion
3. Semantic fallback

---

## Context Rules

Agent MUST:

- include full text only for critical nodes
- summarize neighbors
- avoid loading entire files unnecessarily

---

## Planning Patterns

### Interface-first change

1. update interface
2. update implementations (parallel)
3. update callers
4. remove legacy

---

### Bottom-up change

1. update leaf nodes
2. propagate upward
3. validate

---

### Boundary-inward

1. define boundary
2. work inward toward core logic

---

## Output Format (MANDATORY)

Every task MUST produce:

- Goal
- Affected nodes
- Dependency summary
- SCC groups
- Execution order
- Parallel groups
- Validation plan
- Risks

---

## Final Rule

> The agent does not edit code.
> The agent edits a dependency graph and applies changes through it.
