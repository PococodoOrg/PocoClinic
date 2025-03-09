# Architecture Decision Records (ADR)

This directory contains Architecture Decision Records for the PocoClinic EMR project.

## What is an ADR?

An Architecture Decision Record (ADR) is a document that captures an important architectural decision made along with its context and consequences.

ADRs help team members and stakeholders understand:
- Why a particular decision was made
- What alternatives were considered
- What trade-offs were accepted
- What context existed at the time

## ADR Format

Each ADR follows this format:

```markdown
# ADR-NNNN: Title

## Status
[Proposed, Accepted, Deprecated, Superseded]

## Context
What is the issue that we're seeing that is motivating this decision or change?

## Decision
What is the change that we're proposing and/or doing?

## Consequences
What becomes easier or more difficult to do because of this change?
```

## File Naming

ADRs are numbered sequentially and named using the format:
`NNNN-title-with-hyphens.md`

Example: `0001-modular-monolith-architecture.md`

## Viewing ADRs

All ADR files in this directory are markdown files that can be read directly in your code editor or on GitHub. 