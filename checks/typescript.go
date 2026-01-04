package checks

// TypeScriptChecker implements checks for TypeScript/JavaScript projects.
type TypeScriptChecker struct{}

// Name returns the checker name.
func (c *TypeScriptChecker) Name() string {
	return "TypeScript"
}

// Check runs TypeScript checks on the specified directory.
func (c *TypeScriptChecker) Check(dir string, opts Options) []Result {
	var results []Result

	// Determine package manager
	pm := c.detectPackageManager(dir)

	// Lint check
	if opts.Lint {
		results = append(results, c.checkLint(dir, pm))
	}

	// Format check (prettier)
	if opts.Format {
		results = append(results, c.checkFormat(dir, pm))
	}

	// Type check
	results = append(results, c.checkTypes(dir, pm))

	// Test check
	if opts.Test {
		results = append(results, c.checkTest(dir, pm))
	}

	return results
}

func (c *TypeScriptChecker) detectPackageManager(dir string) string {
	if FileExists(dir + "/pnpm-lock.yaml") {
		return "pnpm"
	}
	if FileExists(dir + "/yarn.lock") {
		return "yarn"
	}
	if FileExists(dir + "/bun.lockb") {
		return "bun"
	}
	return "npm"
}

func (c *TypeScriptChecker) checkLint(dir string, pm string) Result {
	name := "TypeScript: eslint"

	// Try running lint script from package.json first
	result := RunCommand(name, dir, pm, "run", "lint", "--if-present")
	if result.Error == nil {
		return result
	}

	// Fall back to direct eslint
	if CommandExists("eslint") {
		return RunCommand(name, dir, "eslint", ".")
	}

	return Result{
		Name:    name,
		Skipped: true,
		Reason:  "eslint not configured",
	}
}

func (c *TypeScriptChecker) checkFormat(dir string, pm string) Result {
	name := "TypeScript: prettier"

	// Try running format:check script from package.json
	result := RunCommand(name, dir, pm, "run", "format:check", "--if-present")
	if result.Error == nil {
		return result
	}

	// Try prettier:check
	result = RunCommand(name, dir, pm, "run", "prettier:check", "--if-present")
	if result.Error == nil {
		return result
	}

	// Fall back to direct prettier
	if CommandExists("prettier") {
		return RunCommand(name, dir, "prettier", "--check", ".")
	}

	return Result{
		Name:    name,
		Skipped: true,
		Reason:  "prettier not configured",
	}
}

func (c *TypeScriptChecker) checkTypes(dir string, pm string) Result {
	name := "TypeScript: type check"

	// Check if tsconfig.json exists
	if !FileExists(dir + "/tsconfig.json") {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "no tsconfig.json found",
		}
	}

	// Try running typecheck script
	result := RunCommand(name, dir, pm, "run", "typecheck", "--if-present")
	if result.Error == nil {
		return result
	}

	// Try tsc --noEmit
	return RunCommand(name, dir, "npx", "tsc", "--noEmit")
}

func (c *TypeScriptChecker) checkTest(dir string, pm string) Result {
	name := "TypeScript: tests"

	// Try running test script
	result := RunCommand(name, dir, pm, "run", "test", "--if-present")
	if result.Error != nil && result.Output == "" {
		return Result{
			Name:    name,
			Skipped: true,
			Reason:  "no test script configured",
		}
	}

	return result
}
