---
description: Run dev-doctor diagnostics and intelligently fix issues
---

You are helping the user run dev-doctor diagnostics and fix development environment issues.

## Your Task

Run dev-doctor in an intelligent, interactive way that goes beyond the automated cures when needed.

## Step-by-Step Process

### 1. Ask User for Profile

Ask the user which profile they want to check:
- **basic** - Core developer tools (Git, AWS, VPN)
- **infrastructure** - Platform engineering tools (Docker, OpenTofu, Brewfile)
- **data** - Data engineering tools (Python, Wasp)

### 2. Build dev-doctor

Before running, build the latest version:
```bash
cd /Users/luka.borec/PycharmProjects/dev-doctor
go build -o dev-doctor ./cmd/dev-doctor
```

### 3. Run Diagnostics Only First

Run dev-doctor in diagnosis mode to see all issues:
```bash
./dev-doctor --profile <PROFILE> --mode diagnosis
```

Analyze the output and identify which diagnostics failed.

### 4. Apply Cures Systematically

For EACH failing diagnostic (in order):

a. **Run the cure:**
```bash
./dev-doctor --profile <PROFILE> --mode treatment
```

Note: This will run ALL cures, but you should monitor each one.

b. **Check if cure succeeded:**
- Run diagnosis again to verify the specific check now passes
- If it passes: ✓ Move to next failing diagnostic
- If it fails: Continue to step c

c. **Intelligent problem-solving:**
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

6. **Document what you did:**
   - Explain to the user what the problem was
   - Explain how you fixed it
   - Suggest if the cure code should be updated

### 5. Final Summary

After processing all failing diagnostics, provide:
- Count of issues fixed automatically
- Count of issues you had to solve manually
- Remaining issues (if any)
- Suggestions for improving the cure code (if applicable)

## Important Rules

1. **One diagnostic at a time** - Don't try to fix everything at once
2. **Verify each fix** - Always re-run diagnostics after applying a cure
3. **Don't give up** - If the automated cure fails, try alternative solutions
4. **Be creative** - You have full system access, use it to solve problems
5. **Learn from failures** - Suggest code improvements to prevent future issues
6. **User communication** - Keep user informed of progress and what you're trying

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
- Suggestion: Update start_docker cure to actually start Docker instead of just showing instructions
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
