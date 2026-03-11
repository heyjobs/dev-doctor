---
description: Interactive dev-doctor diagnostics with Claude-assisted fixing
---

You are helping the user run dev-doctor diagnostics and interactively fix specific issues they choose.

## Workflow

### 1. Ask User for Profile

Ask the user which profile they want to check:
- **basic** - Core developer tools (Homebrew, AWS, VPN, Claude CLI)
- **infrastructure** - Platform tools (Docker, docker-compose, OpenTofu) + basic
- **data** - Data engineering tools (Python, Wasp) + basic

### 2. Run Diagnostics

Run dev-doctor in diagnosis mode from the project directory:
```bash
dev-doctor --profile <PROFILE> --mode diagnosis
```

Show the user the complete output with all diagnostic results.

### 3. Ask Which Issue to Fix

Look at the diagnostics results and identify issues that are WARNING or CRITICAL (ignore INFO).

Ask the user: **"Which issue would you like me to help fix?"**

Present the failing diagnostics as options for the user to choose from.

### 4. Fix the Chosen Issue

Once the user selects an issue:

a. **Try the automated cure first** (if available):
   ```bash
   dev-doctor --profile <PROFILE> --mode treatment
   ```
   This will run the automated cure for that issue.

b. **Verify if it worked:**
   ```bash
   dev-doctor --profile <PROFILE> --mode diagnosis
   ```
   Check if the issue is now resolved.

c. **If automated cure worked** ✓
   - Confirm with user
   - Ask if they want to fix another issue

d. **If automated cure failed** - Use intelligent problem-solving:

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

   4. **Research if needed:**
      - Use web search to find known issues
      - Check official documentation
      - Look for release notes about compatibility

   5. **Execute the fix:**
      - Try your alternative solution
      - Verify it works with another diagnostic run

   6. **Explain what you did:**
      - Tell the user what the problem was
      - Explain how you fixed it
      - Document the solution

e. **Ask about next issue:**
   After fixing one issue, ask: "Would you like me to fix another issue?"
   - If yes, show remaining issues and let them choose
   - If no, summarize what was accomplished

### 5. Key Points

- **User chooses which issue to fix** - Don't automatically fix all issues
- **One issue at a time** - Focus on one problem, solve it completely
- **Show all command output** - User needs to see what's happening
- **Beyond automated cures** - When cures fail, use creative problem-solving

## Important Rules

1. **Show terminal output** - Display actual command output so user sees what's happening
2. **User chooses** - Let user select which issue to fix, don't fix everything automatically
3. **One at a time** - Fix one issue completely before moving to the next
4. **Fix environments, not code** - Fix the user's machine, don't modify dev-doctor code
5. **Verify each fix** - Always re-run diagnostics after attempting a fix
6. **Be creative** - When automated cures fail, use alternative solutions

## Example Flow

```
User: /dev-doctor

Claude: I'll help you run dev-doctor diagnostics and fix any issues.

Which profile would you like me to check?
- basic (Homebrew, AWS, VPN, Claude CLI)
- infrastructure (Docker, docker-compose, OpenTofu + basic)
- data (Python, Wasp + basic)

User: infrastructure

Claude: Running diagnostics for infrastructure profile...

[Shows output]
✔ Verify Homebrew is installed                  [HEALTHY]
✖ Verify Docker is installed                    [CRITICAL]
  └─ Docker is not installed
ℹ Verify Docker daemon is running               [INFO]
⚠ Verify docker-compose is installed            [WARNING]
  └─ docker-compose is not installed

I found 2 issues that need attention:
1. Docker is not installed [CRITICAL]
2. docker-compose is not installed [WARNING]

Which issue would you like me to help fix?

User: Docker

Claude: I'll help fix the Docker installation issue.

First, let me try the automated cure...
[Runs: dev-doctor --profile infrastructure --mode treatment]
[Shows output]

The automated cure failed because your macOS version is too old.
Let me fix this manually...

[Downloads compatible Docker version]
[Installs from DMG]

✓ Docker Desktop installed successfully

Let me verify it worked...
[Runs: dev-doctor --profile infrastructure --mode diagnosis]

✓ Docker is now installed and working!

Would you like me to fix another issue? (You still have docker-compose not installed)

User: Yes

Claude: I'll fix the docker-compose issue...
[Continues with next fix]
```

## Context

Check `configs/diagnostics.yaml` to understand:
- Which diagnostics exist
- Which cures are available
- Severity levels
- Profile groupings
