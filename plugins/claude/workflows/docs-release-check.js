// Canonical source: plexusone/agent-team-release .claude/workflows/docs-release-check.js
// Mirrored byte-for-byte at grokify/grokify-agent-config claude-code/workflows/docs-release-check.js
// (symlinked into ~/.claude/workflows/ via `agentcfg install claude-code`, which is what
// actually executes). Edit here first — this is the reusable, org-wide version — then copy
// the change to the mirror. Also duplicated at plugins/claude/workflows/ for plugin
// distribution; keep all three in sync.

export const meta = {
  name: 'docs-release-check',
  description: 'Review and fix documentation for release readiness, loop until clean',
  whenToUse: 'Before tagging a release, to ensure all docs are updated',
  phases: [
    { title: 'Review', detail: 'docs-reviewer finds missing/outdated docs' },
    { title: 'Fix', detail: 'docs-writer fixes identified issues' },
    { title: 'Verify', detail: 'docs-reviewer confirms fixes' },
  ],
}

// Args:
//   repoPath: string - path to repository (default: ".")
//   version: string - REQUIRED. Version being released (e.g., "v0.1.0").
//     Resolve this before calling: `git tag --sort=-v:refname | head -1` for the
//     previous tag, then `schangelog parse-commits --since=<tag>` to pick the next
//     version (feat→minor, fix→patch, breaking→major). Confirm with the user.
//   previousVersion: string - last released version tag (e.g., "v0.0.1")
//     If omitted, agents auto-detect from `git tag --sort=-v:refname`.
//     If no tags exist (first release), agents use all commits.
//   features: string - newline-separated list of new features to document
//   maxIterations: number - max review/fix cycles (default: 3)

// The harness can deliver `args` as a JSON-encoded STRING rather than a parsed
// object — observed 2026-09-07: an object passed to the Workflow tool arrived
// here as the string '{"version": "v0.21.0"}', so `args?.version` was always
// undefined and the guard below rejected every invocation. Normalize the three
// shapes: parsed object, JSON-string object, bare version string ("v1.2.0").
let wfArgs = args
if (typeof wfArgs === 'string') {
  const s = wfArgs.trim()
  if (s.startsWith('{')) {
    try { wfArgs = JSON.parse(s) } catch { wfArgs = {} }
  } else {
    wfArgs = /^v?\d+\.\d+\.\d+$/.test(s) ? { version: s } : {}
  }
}

if (!wfArgs?.version) {
  throw new Error(
    "docs-release-check requires args.version (e.g. 'v1.2.0'; a bare version string " +
    "as the whole args value is also accepted). " +
    "It used to default to the literal string 'unreleased' when omitted, which made every " +
    "check vacuous — the review had nothing real to validate against, so it always returned " +
    "ready:true on the first pass without creating a CHANGELOG entry, release notes, or " +
    "mkdocs nav update. Resolve the version first (git tags + `schangelog parse-commits`, " +
    "or ask the user) and pass it explicitly."
  )
}

const MAX_ITERATIONS = wfArgs?.maxIterations || 3
const repoPath = wfArgs?.repoPath || '.'
const version = wfArgs.version
const previousVersion = wfArgs?.previousVersion || ''
const features = wfArgs?.features || ''

const repoPathGuard = `Before doing anything else, run: test -d "${repoPath}" && cd "${repoPath}" && pwd
If that directory does not exist, or you cannot cd into it, STOP and report it as a blocking
issue named "repoPath" — do not fall back to reviewing or editing any other directory
(including your own ambient working directory). All file paths below are relative to
${repoPath}.`

const commitRangeHint = previousVersion
  ? `Previous version: ${previousVersion}
Use: schangelog parse-commits --since=${previousVersion}`
  : `No previous version specified. Detect the latest tag with:
  git tag --sort=-v:refname | head -1
If no tags exist (first release), use:
  schangelog parse-commits --last=100
or parse all commits since the root.`

const changelogJsonSpec = `CHANGELOG.json format (structured-changelog):
{
  "irVersion": "1.0",
  "project": "<repo-name>",
  "repository": "<repo-url>",
  "versioning": "semver",
  "commitConvention": "conventional",
  "releases": [{
    "version": "${version}",
    "date": "<YYYY-MM-DD>",
    "highlights": [{"description": "..."}],
    "added": [{"description": "...", "commit": "<hash>"}],
    "fixed": [...], "changed": [...], "refactored": [...],
    "documentation": [...], "dependencies": [...], "tests": [...]
  }]
}
Derive "project" and "repository" from: git remote get-url origin
  e.g. https://github.com/org/repo → project: "repo", repository: "https://github.com/org/repo"
  The local directory name may differ from the GitHub repo name — always use the remote.
Category mapping: feat→added, fix→fixed, refactor→refactored,
docs→documentation, chore(deps)→dependencies, test→tests, perf→changed.
Commit hashes are required for each entry.`

const ISSUES_SCHEMA = {
  type: 'object',
  properties: {
    ready: { type: 'boolean', description: 'True if docs are release-ready' },
    issues: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          file: { type: 'string', description: 'File path relative to repo' },
          severity: { type: 'string', enum: ['blocking', 'should-fix', 'nice-to-have'] },
          description: { type: 'string', description: 'What needs to be fixed' },
        },
        required: ['file', 'severity', 'description']
      }
    }
  },
  required: ['ready', 'issues']
}

const FIX_SCHEMA = {
  type: 'object',
  properties: {
    fixed: { type: 'array', items: { type: 'string' }, description: 'Files that were fixed' },
    skipped: { type: 'array', items: { type: 'string' }, description: 'Files skipped with reason' },
  },
  required: ['fixed', 'skipped']
}

let iteration = 0
let lastReview = null

while (iteration < MAX_ITERATIONS) {
  iteration++
  phase('Review')
  log(`Iteration ${iteration}/${MAX_ITERATIONS}: Running docs review`)

  const reviewPrompt = iteration === 1
    ? `Review documentation in ${repoPath} for ${version} release readiness.

${repoPathGuard}

${commitRangeHint}

New features to document:
${features}

Check: README, CHANGELOG.json, CHANGELOG.md, release notes, mkdocs.yml nav.

For CHANGELOG validation:
- If CHANGELOG.json exists, run: schangelog validate CHANGELOG.json
- If CHANGELOG.json does not exist, report it as blocking — it must be created.
- If CHANGELOG.md is missing or doesn't include ${version}, report as blocking.
- Release notes: check docs/releases/${version}.md (if docs/ exists) or RELEASE_NOTES_${version}.md

Report blocking issues that must be fixed before tagging ${version}.`
    : `Re-review documentation in ${repoPath} after fixes.

${repoPathGuard}

Previous issues found: ${JSON.stringify(lastReview?.issues || [])}

Verify the fixes were applied correctly. Check:
- schangelog validate CHANGELOG.json passes
- CHANGELOG.md includes ${version} and was regenerated from JSON
- Release notes exist and are accurate
- mkdocs.yml nav is updated (if applicable)

Report any remaining issues.`

  const review = await agent(reviewPrompt, {
    agentType: 'docs-reviewer',
    label: `review-${iteration}`,
    phase: 'Review',
    schema: ISSUES_SCHEMA,
  })

  if (!review) {
    log('Review agent failed')
    break
  }

  lastReview = review

  if (review.ready) {
    log(`Documentation is release-ready after ${iteration} iteration(s)`)
    return { ready: true, iterations: iteration, review }
  }

  const blockingIssues = review.issues.filter(i => i.severity === 'blocking')
  if (blockingIssues.length === 0) {
    log('No blocking issues remain')
    return { ready: true, iterations: iteration, review }
  }

  log(`Found ${blockingIssues.length} blocking issues to fix`)

  phase('Fix')
  const fixPrompt = `Fix documentation issues in ${repoPath} for ${version} release.

${repoPathGuard}

${commitRangeHint}

Issues to fix (blocking only):
${JSON.stringify(blockingIssues, null, 2)}

New features that need documentation:
${features}

Fix each blocking issue. Create missing files, update existing docs.

IMPORTANT: Derive the repository URL from git, not the local directory name:
  git remote get-url origin
The local directory may differ from the GitHub repo name.

For CHANGELOG.json:
- Parse commits first to get hashes and types.
- If creating from scratch, use the full structured-changelog format:
${changelogJsonSpec}
- If updating, add a new release entry at the top of the releases array.

For CHANGELOG.md:
- After creating/updating CHANGELOG.json, regenerate with:
  schangelog generate CHANGELOG.json -o CHANGELOG.md
- If schangelog is not available, write CHANGELOG.md manually matching the JSON content.

For release notes:
- Create at docs/releases/${version}.md (if docs/ exists) or RELEASE_NOTES_${version}.md
- Include: summary, highlights, what's new, installation command.

For mkdocs.yml:
- Add ${version} to the Releases nav section (if mkdocs.yml exists).

Follow existing doc style. Do not modify code files.`

  const fixes = await agent(fixPrompt, {
    agentType: 'docs-writer',
    label: `fix-${iteration}`,
    phase: 'Fix',
    schema: FIX_SCHEMA,
  })

  if (!fixes || fixes.fixed.length === 0) {
    log('No fixes applied, stopping loop')
    break
  }

  log(`Fixed ${fixes.fixed.length} files: ${fixes.fixed.join(', ')}`)
}

log(`Reached max iterations (${MAX_ITERATIONS})`)
return { ready: false, iterations: iteration, review: lastReview }
