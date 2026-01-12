# Release Team Orchestration

Release validation and deployment workflow

**Process:** hierarchical
**Manager:** release-coordinator

---

## Group 1 (Sequential)

PM validation runs first to determine the target version.

### Task: pm-validation

Product scope and version recommendation

**Instructions:**

Use the Task tool to spawn subagent `pm`:

```
Task tool:
  subagent_type: general-purpose
  description: "Product scope and version recommendation"
  prompt: |
    You are the pm specialist. Read your instructions from validation/specs/pm.md

    Execute the following subtasks and report Go/No-Go for each:
    - version-recommendation
    - release-scope
    - changelog-quality
    - breaking-changes
    - roadmap-alignment
    - deprecation-notices

    Output the recommended version number (e.g., v0.3.0) for use by downstream tasks.
```

**Subtasks:**

| Subtask | Type | Required | Expected |
|---------|------|----------|----------|
| version-recommendation | command | Yes | Recommended version determined |
| release-scope | file | No | Scope aligned with roadmap |
| changelog-quality | file | Yes | Entries are user-facing and... |
| breaking-changes | command | Yes | All breaking changes docume... |
| roadmap-alignment | file | No | Release aligns with roadmap |
| deprecation-notices | pattern | No | Deprecations documented |

**Sign-off:** GO if all 3 required subtasks pass. Optional subtasks report WARN on failure.

**Output:** Confirmed version number (e.g., `v0.3.0`) passed to downstream tasks.

---

## Parallel Group 2

These tasks can run concurrently using parallel Task tool calls. They depend on pm-validation completing first.

**Requires:** pm-validation (must be GO)

### Task: qa-validation

Quality Assurance validation

**Instructions:**

Use the Task tool to spawn subagent `qa`:

```
Task tool:
  subagent_type: general-purpose
  description: "Quality Assurance validation"
  prompt: |
    You are the qa specialist. Read your instructions from validation/specs/qa.md

    Execute the following subtasks and report Go/No-Go for each:
    - build
    - tests
    - lint
    - format
    - mod-tidy
    - error-handling
    - local-replace

    Target version: {VERSION from pm-validation}
```

**Subtasks:**

| Subtask | Type | Required | Expected |
|---------|------|----------|----------|
| build | command | Yes | No errors |
| tests | command | Yes | All tests pass |
| lint | command | Yes | No lint errors |
| format | command | Yes | No output (all files forma... |
| mod-tidy | command | Yes | No diff output |
| error-handling | pattern | Yes | No matches (errors should... |
| local-replace | pattern | Yes | No matches |

**Sign-off:** GO if all 7 required subtasks pass. Optional subtasks report WARN on failure.

### Task: docs-validation

Documentation validation

**Instructions:**

Use the Task tool to spawn subagent `documentation`:

```
Task tool:
  subagent_type: general-purpose
  description: "Documentation validation"
  prompt: |
    You are the documentation specialist. Read your instructions from validation/specs/documentation.md

    Execute the following subtasks and report Go/No-Go for each:
    - readme
    - changelog-md
    - changelog-json
    - release-notes
    - prd
    - trd

    Target version: {VERSION from pm-validation}
```

**Subtasks:**

| Subtask | Type | Required | Expected |
|---------|------|----------|----------|
| readme | file | Yes | Pass |
| changelog-md | file | Yes | Pass |
| changelog-json | file | No | Pass |
| release-notes | pattern | Yes | Pass |
| prd | file | No | Pass |
| trd | file | No | Pass |

**Sign-off:** GO if all 3 required subtasks pass. Optional subtasks report WARN on failure.

---

## Execution Flow

```
═══════════════════════════════════════════════════════════════
                        PHASE 1: Review
═══════════════════════════════════════════════════════════════

                    ┌─────────────────┐
                    │  pm-validation  │
                    │   (version)     │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
       ┌──────────┐   ┌──────────┐   ┌──────────┐
       │    qa    │   │   docs   │   │ security │
       └────┬─────┘   └────┬─────┘   └────┬─────┘
            │              │              │
            └──────────────┼──────────────┘
                           ▼
                    ┌─────────────────┐
                    │ release-valid.  │
                    └─────────────────┘

        ▼▼▼ Fix issues and commit changes ▼▼▼

═══════════════════════════════════════════════════════════════
                     PHASE 2: Finalize
═══════════════════════════════════════════════════════════════

                    ┌─────────────────┐
                    │changelog-finalize│
                    │ (link commits)  │
                    └────────┬────────┘
                             ▼
                    ┌─────────────────┐
                    │ execute-release │
                    │  (tag & push)   │
                    └─────────────────┘
```

### Human-in-the-Loop Prompts

During Phase 1, agents may require human input:

| Agent | Trigger | Prompt |
|-------|---------|--------|
| PM | CHANGELOG.json missing target version | Create entry (auto-generate / manual / abort) |
| PM | Uncommitted changes only | Confirm version recommendation |
| Docs | Release notes missing (major/minor) | Create release notes at correct path |

**CHANGELOG.json creation workflow:**

```
PM detects missing v0.2.0 entry
         │
         ▼
  ┌─────────────────┐
  │  HITL Prompt:   │
  │  1. Auto-generate│
  │  2. Manual entry │
  │  3. Abort        │
  └────────┬────────┘
           │
     ┌─────┴─────┐
     ▼           ▼
Auto-generate  Manual
     │           │
     ▼           ▼
Draft entry   Human types
     │        entries
     ▼           │
Human reviews ◄──┘
     │
     ▼
Write CHANGELOG.json
     │
     ▼
Continue validation
```

### Phase 1 → Phase 2 Transition

After Phase 1 validations pass:

1. Fix any issues identified by agents
2. **Break changes into topical commits** (see release.md for commit conventions)
3. Commit all fixes with conventional commit messages
4. `changelog-finalize` populates commit hashes in CHANGELOG.json
5. `changelog-finalize` generates CHANGELOG.md
6. `changelog-finalize` commits changelog updates
7. `execute-release` creates and pushes the tag

---

## Expected Status Report

**Report width:** 78 characters (fits 80-column terminals)

**Status message guidelines:**
- Check name column: 22 characters max
- Status icon + label: 8 characters (e.g., `🟢 GO`)
- Message column: ~35 characters max

After execution, report status in this format:

```
╔════════════════════════════════════════════════════════════════════════════╗
║                            TEAM STATUS REPORT                              ║
╠════════════════════════════════════════════════════════════════════════════╣
║ Project: github.com/grokify/release-agent                                  ║
║ Target:  v0.3.0                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ PHASE 1: REVIEW                                                            ║
╠════════════════════════════════════════════════════════════════════════════╣
║ pm-validation (pm)                                                         ║
║   version-recommendation   🟢 GO    v0.3.0 (minor bump)                    ║
║   release-scope            🟢 GO    (optional)                             ║
║   changelog-quality        🟢 GO    Highlights present                     ║
║   breaking-changes         🟢 GO    None detected                          ║
║   roadmap-alignment        🟢 GO    (optional)                             ║
║   deprecation-notices      🟢 GO    (optional)                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ qa-validation (qa)                                                         ║
║   build                    🟢 GO                                           ║
║   tests                    🟢 GO                                           ║
║   lint                     🟢 GO                                           ║
║   format                   🟢 GO                                           ║
║   mod-tidy                 🟢 GO                                           ║
║   error-handling           🟢 GO                                           ║
║   local-replace            🟢 GO                                           ║
╠════════════════════════════════════════════════════════════════════════════╣
║ docs-validation (documentation)                                            ║
║   readme                   🟢 GO    340 lines, comprehensive               ║
║   changelog-md             🟢 GO    v0.3.0 documented                      ║
║   changelog-json           🟢 GO    (optional)                             ║
║   release-notes            🟢 GO    docs/releases/v0.3.0.md                ║
║   prd                      🟡 SKIP  (optional)                             ║
║   trd                      🟡 SKIP  (optional)                             ║
╠════════════════════════════════════════════════════════════════════════════╣
║ PHASE 2: FINALIZE                                                          ║
╠════════════════════════════════════════════════════════════════════════════╣
║ changelog-finalize (release)                                               ║
║   link-commits             🟢 GO    12 commits linked                      ║
║   generate-changelog       🟢 GO    CHANGELOG.md updated                   ║
║   commit-changelog         🟢 GO    Committed abc1234                      ║
╠════════════════════════════════════════════════════════════════════════════╣
║ execute-release (release-coordinator)                                      ║
║   ci-status                🟢 GO    CI passed                              ║
║   create-tag               🟢 GO    v0.3.0 created                         ║
║   push-tag                 🟢 GO    Pushed to origin                       ║
╠════════════════════════════════════════════════════════════════════════════╣
║                          🚀 TEAM: GO for v0.3.0 🚀                         ║
╚════════════════════════════════════════════════════════════════════════════╝
```
