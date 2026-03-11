---
description: Manually run dev-doctor diagnostics and intelligently fix issues
---

**Note**: This is for MANUAL/INTERACTIVE use. In production, dev-doctor's treatment mode automatically spawns Claude when cures fail (see CLAUDE.md for details).

You are helping the user manually run dev-doctor diagnostics and fix development environment issues.

## Your Task

Run dev-doctor in an intelligent, interactive way that goes beyond the automated cures when needed.

## Step-by-Step Process

### 1. Ask User for Profile

Ask the user which profile they want to check:
- **basic** - Core developer tools (Git, AWS, VPN)
- **infrastructure** - Platform engineering tools (Docker, OpenTofu, Brewfile)
- **data** - Data engineering tools (Python, Wasp)

### 2. Build dev-doctor

Before running, build the latest version from the dev-doctor repository root:
```bash
go build -o dev-doctor ./cmd/dev-doctor
```

### 3. Run Diagnostics to Identify Issues

Run dev-doctor in diagnosis mode and **respond with the full output**:
```bash
./dev-doctor --profile <PROFILE> --mode diagnosis
```

Show the user the complete terminal output. Identify which diagnostics are WARNING or CRITICAL.

### 4. Fix Issues One at a Time

**IMPORTANT:** Process diagnostic-cure pairs ONE AT A TIME, not all at once.

For EACH failing diagnostic (WARNING or CRITICAL):

a. **Show the diagnostic result** - Tell user which issue you're fixing

b. **Run only diagnostics to see current status:**
```bash
./dev-doctor --profile <PROFILE> --mode diagnosis
```
Show the relevant output for this specific diagnostic.

c. **Run treatment mode** (this will apply all available cures):
```bash
./dev-doctor --profile <PROFILE> --mode treatment
```
**Show the full terminal output** - user needs to see what the cure is doing.

d. **Verify it worked** - Run diagnostics again:
```bash
./dev-doctor --profile <PROFILE> --mode diagnosis
```
Show the output and confirm if this specific diagnostic now passes.

e. **If cure succeeded** ✓ - Move to next failing diagnostic

f. **If cure failed** - Continue to step 4g

g. **Intelligent problem-solving:**
When a cure fails, you must go beyond the scripted cure:

1. **Analyze the error:**
   - Read error messages carefully
   - Identify root cause (permissions, network, version incompatibility, etc.)

2. **Check prerequisites:**
   - Is Homebrew installed and working?
   - Are there permission issues?
   - Is the network accessible?
   - Is the macOS version compatible?

3. **Try alternative solutions:**
   - Manual installation via DMG/PKG
   - Different installation methods (brew vs direct download)
   - Older/newer versions that might work
   - Configuration fixes

4. **Research the issue:**
   - Use web search to find known issues
   - Check official documentation
   - Look for release notes about compatibility

5. **Execute the fix:**
   - Try your alternative solution
   - Verify it works with another diagnostic run

6. **Explain what you did:**
   - Tell the user what the problem was
   - Explain how you fixed it
   - Make sure they understand the solution

### 5. Final Summary

After processing all failing diagnostics, provide:
- Count of issues fixed automatically
- Count of issues you had to solve manually
- Remaining issues (if any)
- Brief explanation of manual fixes applied

## Important Rules

1. **Respond with terminal output** - User needs to see the actual program output from each diagnostic and cure run
2. **One diagnostic at a time** - Process each failing diagnostic-cure pair separately, don't batch them
3. **Fix environments, not code** - Don't suggest changes to dev-doctor's cure code. Cures work in most cases, but environments differ and sometimes need manual intervention.
4. **Verify each fix** - Always re-run diagnostics after applying a cure
5. **Don't give up** - If the automated cure fails, try alternative solutions
6. **Be creative** - You have full system access, use it to solve problems

## Example Flow

```
User: Run dev-doctor for infrastructure profile

You:
1. Building dev-doctor...
   ✓ Built successfully

2. Running diagnostics...
   Found 3 failing checks:
   - Docker not installed [CRITICAL]
   - Docker not running [INFO]
   - docker-compose not installed [WARNING]

3. Fixing: Docker not installed
   - Attempting automated cure (install_docker)...
   - Homebrew installation failed (macOS 13 too old)
   - Detected macOS 13.4 - need compatible Docker version
   - Downloading Docker Desktop 4.43.0 for macOS 13...
   - Installing from DMG...
   ✓ Docker Desktop installed successfully

4. Re-running diagnostics...
   ✓ Docker is now installed
   - Next issue: Docker not running

5. Fixing: Docker not running
   - Cure shows instructions but doesn't auto-start
   - Starting Docker Desktop automatically...
   - Waiting for daemon to be ready...
   ✓ Docker daemon is running

6. Re-running diagnostics...
   ✓ All checks passing!

Summary:
- Fixed 3 issues (2 automatically, 1 with manual intervention)
- All infrastructure checks now passing
- Manual fix: Started Docker daemon automatically instead of just showing instructions
```

## Context

Read the CLAUDE.md file for:
- Project structure
- How diagnostics and cures work
- Status vs severity concepts
- Best practices for cures

Check `configs/diagnostics.yaml` to understand:
- Which diagnostics map to which cures
- Severity levels
- Profile groupings
