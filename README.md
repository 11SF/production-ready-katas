# production-ready-katas

A collection of hands-on coding exercises for engineers who already know how to code — but want to understand what's *really* happening when their code runs in production.

---

## Who this is for

Senior engineers who are fluent in one language and picking up another (e.g. Go, Rust, Zig), but find that tutorials and AI-generated code leave gaps in their mental model of:

- What the OS is actually doing when you open a file
- Why the naive approach works locally but breaks in production
- How to read stdlib documentation instead of copying examples

If you've ever shipped code that passed all local tests but caused a `too many open files` at 2am — this is for you.

---

## Philosophy

**1. See the problem before the solution**

Every kata starts with the naive approach and why it breaks. You don't get the "right way" until you understand the failure mode.

**2. Real incidents, not toy examples**

Each exercise is grounded in real production incidents — OOM loops, fd leaks, symlink attacks. The context is always "this actually happened."

**3. Read the source, not the tutorial**

The `Explore First` section in each kata gives you method names to look up in the stdlib — not examples to copy. You're expected to open `go to definition` and read.

**4. One spec, multiple languages**

Problem definitions live in `shared/` and are language-agnostic. Each language has its own implementation folder. Same kata, different idioms — good for comparison.

---

## Structure

```
production-ready-katas/
├── shared/
│   ├── concepts/          # Background reading: OS behavior, memory, I/O
│   ├── assignment-specs/  # Language-agnostic kata definitions
│   │   ├── file-handling/
│   │   ├── compression-encryption/
│   │   ├── networking/
│   │   ├── cloud-storage/
│   │   └── concurrency/
│   └── scenarios/         # End-to-end exercises combining multiple katas
├── go/                    # Go implementations
├── rust/                  # Rust implementations
└── zig/                   # Zig implementations
```

### Kata spec format

Each kata follows this structure:

| Section | Purpose |
|---|---|
| **Context** | Why this problem matters in real systems |
| **Real World Incidents** | Actual production failures caused by this mistake |
| **The Naive Way** | What most people write first, and exactly where it breaks |
| **Explore First** | Method names to look up in stdlib — no examples, just hints |
| **Task** | What you need to build |
| **Requirements** | Constraints that prevent you from taking shortcuts |
| **Acceptance Criteria** | Testable checklist — all must pass |
| **Concepts Involved** | Links to background reading in `shared/concepts/` |

---

## Domains

Katas are grouped by domain. Numbering within each group goes from naive baseline (01) to increasingly complex patterns — not overall difficulty.

| Domain | Description |
|---|---|
| `file-handling` | Read patterns, write patterns, encoding, resource management, edge cases, integrity, testing |
| `compression-encryption` | Stream-based compress/decompress, encrypt/decrypt for large files |
| `networking` | Download with retry/resume, upload/multipart, pipe without temp files |
| `cloud-storage` | S3/GCS-style object storage patterns |
| `concurrency` | File locking, TOCTOU race conditions |

---

## Start here

Not sure where to begin? Read these in order:

1. **[Whole-File Read kata](shared/assignment-specs/file-handling/01-read-patterns/01-whole-file-read.md)** — the first kata, a good example of the full format
2. **[fd-lifecycle](shared/concepts/fd-lifecycle.md)** — background reading referenced by the kata above
3. Pick a language folder (`go/`, `rust/`, `zig/`) and implement it

The kata spec tells you what to build. `Explore First` tells you where to look. The concept docs explain the *why*.

---

## Current content

### Katas available

| Kata | Domain | Difficulty |
|---|---|---|
| [Whole-File Read](shared/assignment-specs/file-handling/01-read-patterns/01-whole-file-read.md) | file-handling / read-patterns | 1 |

### Concept docs available

| Concept | Description |
|---|---|
| [fd-lifecycle](shared/concepts/fd-lifecycle.md) | File descriptors, OS limits, `/proc`, fork inheritance, strace |
| [error-wrapping](shared/concepts/error-wrapping.md) | `%w`, `errors.Is/As`, errno, syscall error chain |
| [memory-allocation](shared/concepts/memory-allocation.md) | Heap, page cache, virtual memory, OOM Killer, mmap |

---

## Roadmap

### Near term
- [ ] Go starter code + test file for `01-whole-file-read`
- [ ] Kata: `02-streaming-read` — reading large files without loading into memory
- [ ] Kata: `01-basic-write` — write patterns and atomicity
- [ ] Concept doc: bytes, characters, and encoding fundamentals
- [ ] Concept doc: what "stream" actually means (not just files — stdin, network, pipe)

### Medium term
- [ ] GitHub Actions CI — auto-run tests per language on push
- [ ] Rust implementations for file-handling katas
- [ ] First scenario: "process a large file from S3, resume if interrupted"
- [ ] `Explore First` sections split per language (Go / Rust / Zig)

### Open decisions
- Kata template: lean (4 sections) vs full (all sections) — currently using full for all
- Scenario prerequisites: hard gate or soft recommendation?
- TDD starter tests: provide failing tests per kata, or let engineers write their own?

---

## Contributing

If you want to add a kata, concept doc, or language implementation:

1. Kata specs go in `shared/assignment-specs/<domain>/` — keep them language-agnostic
2. Follow the section format above — especially `Real World Incidents` and `Explore First`
3. Concept docs go in `shared/concepts/` — include both approachable explanation and OS-level depth
4. Implementations go in `<language>/<domain>/` mirroring the spec path
