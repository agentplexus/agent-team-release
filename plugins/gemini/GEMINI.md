# Release Agent

This plugin automates software release workflows for multi-language repositories.

## Available Commands

- `/release-agent:release <version>` - Execute full release workflow
- `/release-agent:check` - Run validation checks
- `/release-agent:changelog <version>` - Generate changelog
- `/release-agent:version-next` - Suggest next version

## Supported Languages

- Go (build, test, golangci-lint, gofmt)
- TypeScript/JavaScript (ESLint, Prettier, npm test)

## Dependencies

- `atrelease` - Release automation CLI
- `schangelog` - Changelog generation
- `git` - Version control
- `gh` - GitHub CLI