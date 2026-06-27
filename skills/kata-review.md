---
name: kata-review
description: Generic review template for any kata implementation — read this before reviewing code
---

# Kata Review Guide

## How to use

Tell the AI: "read shared/skills/kata-review.md first, then review [file]"
The AI will review the code against the dimensions below and present findings in an easy-to-read format.

---

## Review Dimensions

### 1. Correctness
- Does the logic match the kata requirements?
- Are all edge cases specified in the kata handled?
- Is the output correct (byte-for-byte or matching the spec)?

### 2. Resource Management
- Are file descriptors, memory, and goroutines released on every path (both success and error)?
- Is `defer` placed correctly — after the error check, not before?
- Any resources that could leak under concurrent usage?

### 3. Error Handling
- Are all error paths covered?
- Do error messages provide enough context (not just "failed")?
- Is error wrapping correct — can `errors.Is` / `errors.As` still work after wrapping?

### 4. Concurrency Safety
- Is there any shared mutable state?
- If so, is it guarded (mutex, channel, atomic)?
- Is the function safe to call from multiple goroutines simultaneously?

### 5. Language Idioms
- Does naming follow the language's conventions?
- Is the stdlib used appropriately, or is there logic that duplicates what stdlib already provides?
- Does anything feel "un-Go-like" / "un-Rust-like"?

### 6. Production Gap
- In what real-world scenarios would this code fail in production?
- Are there assumptions that hold in dev but break in prod (e.g. file size, concurrent access, network conditions)?
- Is this the pattern production systems actually use, or is there a better pattern or library for this use case?

### 7. Kata Quality
- Is the kata spec still accurate and relevant, or has the stdlib / ecosystem moved on?
- Do the Acceptance Criteria cover all important cases, or are there cases worth adding?
- Is any part of the spec misleading or ambiguous — should any questions be revised?

---

## After Review: Reinforce Understanding

After presenting all dimensions, always do both of the following:

### 8. Explain It Back
Pick the 1-2 most important findings from the review and ask the user to explain them in their own words.

The goal is not to quiz — it's to surface gaps between "I followed the fix" and "I actually understand why."

Example questions:
- "Why does `LimitReader` need `+1` instead of just `MAX_FILE_SIZE`?"
- "What would happen if you called `defer file.Close()` before the error check?"
- "Why does `rename` give us atomicity but `write` doesn't?"

If the user explains correctly — confirm and move on.
If the explanation is off — clarify the concept, link to the relevant concept doc in `shared/concepts/` if one exists.

### 9. Pattern Connections
If any finding in this review resembles something from a previous kata, call it out explicitly.

Example:
- "This is the same fd leak pattern from `01-whole-file-read` — `defer` in the wrong place"
- "The size check issue here is the TOCTOU problem — same class of bug as what `LimitReader` was added to fix"

Connecting patterns across katas builds a mental model, not just isolated fixes.

---

## Output Format

Present each dimension like this:

```
## Review: [kata name]

### [Dimension] [✅ / ⚠️ / ❌]
[Short summary — pass or what the issue is]
[If issue — include line number and what to fix]

---
Must fix before prod: [N] items
Fix when you can: [N] items
Kata spec suggestions: [yes / none]
```

**Legend:**
- ✅ Pass — nothing to change
- ⚠️ Should improve — not critical but worth fixing
- ❌ Must fix — will break in production
