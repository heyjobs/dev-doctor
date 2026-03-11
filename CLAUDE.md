# Claude AI Assistant Configuration

This document defines Claude's behavior when using dev-doctor to fix developer environment issues.

## What is dev-doctor?

**dev-doctor** is a diagnostic and remediation tool that checks developer workstation health and automatically fixes common setup issues. It helps developers get their environments working quickly.

## How It Works

dev-doctor follows a "diagnostic + cure" pattern:

1. **Diagnostics** - Check if something is installed/configured correctly
2. **Cures** - Automatically fix the issue if the diagnostic fails

Each diagnostic returns a status:
- `[INFO]` (blue) - Nice to have, but not critical
- `[WARNING]` (yellow) - Should be fixed, but not blocking
- `[CRITICAL]` (red) - Must be fixed for environment to work

## Profiles

Tests are organized by developer role:
- **basic** - Core tools everyone needs (Git, AWS, VPN)
- **infrastructure** - Platform engineers (Docker, OpenTofu, Brewfile)
- **data** - Data engineers (Python, Wasp)

## The dev-doctor Skill

Use `/dev-doctor` to interactively fix environment issues.

### Expected Behavior

The skill should:

1. **Ask for profile** - Which profile to check (basic/infrastructure/data)

2. **Build dev-doctor** - Always build latest version first from the dev-doctor repository root:
   ```bash
   go build -o dev-doctor ./cmd/dev-doctor
   ```

3. **Run diagnostics** - Check what's broken and **show the full terminal output**:
   ```bash
   ./dev-doctor --profile <PROFILE> --mode diagnosis
   ```

4. **Fix issues one at a time** - For EACH warning or critical issue:

   a. **Show which issue you're fixing**

   b. **Run treatment mode** and **show full terminal output**:
   ```bash
   ./dev-doctor --profile <PROFILE> --mode treatment
   ```

   c. **Verify it worked** - Run diagnostics again and show output:
   ```bash
   ./dev-doctor --profile <PROFILE> --mode diagnosis
   ```

   d. **If cure succeeded** ✓ - Move to next issue

   e. **If cure failed** - You must solve it beyond the automated cure:
      - Analyze the error message
      - Check prerequisites (Homebrew installed? Permissions? Network?)
      - Check OS compatibility (macOS version?)
      - Try alternative installation methods
      - Search for known issues
      - Execute manual fix
      - Verify the fix worked

5. **Report results** - Summary of what was fixed

### Important Rules

- **Show all terminal output** - User needs to see the actual program output from each diagnostic and cure run

- **One issue at a time** - Process each diagnostic-cure pair separately, don't batch them

- **Fix environments, not code** - Don't suggest changes to dev-doctor's cure code. The cures work for most cases, but sometimes environments differ and need manual intervention.

- **Always verify** - Re-run diagnostics after each fix to confirm it worked

- **Don't give up** - If the automated cure fails, figure out an alternative solution

### Example Flow

```
User: /dev-doctor

You: Which profile do you want to check?
- basic (Git, AWS, VPN)
- infrastructure (Docker, OpenTofu)
- data (Python, Wasp)

User: infrastructure

You:
1. Building dev-doctor...
   ✓ Built successfully

2. Running diagnostics...
   Found 2 failing checks:
   - Docker not installed [CRITICAL]
   - docker-compose not installed [WARNING]

3. Fixing: Docker not installed
   - Running install_docker cure...
   - Homebrew installation failed (macOS 13 incompatible)
   - Detected macOS 13.4
   - Trying alternative: Downloading Docker Desktop 4.43.0 for macOS 13...
   - Installing from DMG...
   - Prompting for sudo password...
   ✓ Docker Desktop installed

4. Verifying fix...
   - Re-running diagnostics...
   ✓ Docker is now installed

5. Fixing: docker-compose not installed
   - Running install_docker_compose cure...
   ✓ docker-compose installed via Homebrew

6. Final diagnostics...
   ✓ All infrastructure checks passing!

Summary: Fixed 2 issues successfully
```

## Key Files

- Config: `configs/diagnostics.yaml` - see which diagnostics map to which cures
- Build command (from repo root): `go build -o dev-doctor ./cmd/dev-doctor`
- Run (from repo root): `./dev-doctor --profile <PROFILE> --mode <diagnosis|treatment>`
