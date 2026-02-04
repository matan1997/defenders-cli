# Defenders CLI

A command-line tool for Azure DevOps operations - create work items, PRs, run pipelines, and manage PR approvals.

## Installation

```bash
cd defenders-cli
go build -o defenders .

# Move to PATH (optional)
sudo mv defenders /usr/local/bin/
```

## Quick Start

```bash
# Configure the CLI (interactive setup)
defenders conf

# Create a work item
defenders cado --title="My Feature"

# Create a PR from current branch
defenders prme

# Run a pipeline
defenders release run <pipeline-url>

# Approve a PR
defenders pr --approve <pr-url>
```

## Commands

### `conf` - Configuration

Set up your CLI configuration interactively. Stores settings in:
- **Linux/macOS**: `~/.config/defenders/config.json`
- **Windows**: `%APPDATA%\defenders\config.json`

```bash
# Interactive setup
defenders conf

# Show current configuration
defenders conf show

# Show config file path
defenders conf path

# Reset to defaults
defenders conf reset
```

**Configuration values:**
- PAT Token
- Organization URL (default: `https://dev.azure.com/msazure`)
- Project (default: `One`)
- Team (default: `Rome`)
- Area Path
- Assigned To (email)

---

### `cado` - Create ADO Work Item

Create a Feature work item with automatic iteration assignment.

```bash
# Basic usage
defenders cado --title="Implement new feature"

# With parent link
defenders cado --title="My Task" --parent=12345

# Override assigned-to
defenders cado --title="Feature" --assigned-to="user@microsoft.com"
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--title` | (required) Title of the Feature |
| `--parent` | Parent work item ID to link |
| `--assigned-to` | Override assigned-to from config |

---

### `prme` - Create Pull Request

Create a PR from the current branch to the default branch (develop/main/master).

```bash
# Auto-detect everything
defenders prme

# Link to work item
defenders prme -i 12345

# Custom title
defenders prme -t "My PR Title"

# Both
defenders prme -i 12345 -t "Fix bug in auth module"
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-i, --work-item` | Work item ID to link |
| `-t, --title` | Custom PR title (default: branch name) |

---

### `release` - Pipeline Operations

Run and monitor Azure DevOps pipelines.

#### Run a pipeline

```bash
# Run a pipeline
defenders release run https://dev.azure.com/org/project/_build?definitionId=456

# Run with specific PAT token
defenders release run <url> -t <token>
```

#### Monitor and trigger

Wait for a pipeline to complete, then trigger another:

```bash
defenders release monitor-trigger \
  https://dev.azure.com/org/proj/_build/results?buildId=123 \
  https://dev.azure.com/org/proj/_build?definitionId=456

# Custom interval (default: 30 seconds)
defenders release monitor-trigger <wait-url> <trigger-url> -i 60
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-t, --token` | PAT token (overrides config/env) |
| `-i, --interval` | Check interval in seconds (default: 30) |

---

### `pr` - PR Approval Operations

Approve or reset votes on Pull Requests.

```bash
# Approve a PR
defenders pr --approve https://dev.azure.com/org/project/_git/repo/pullrequest/123

# Reset your vote
defenders pr --reset <pr-url>

# Use another user's PAT (approve on their behalf)
defenders pr --approve <pr-url> -t <other-user-pat>
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--approve` | Approve the PR |
| `--reset` | Reset your vote |
| `-t, --token` | PAT token (use another user's PAT) |

---

## Authentication

### Option 1: Configuration file (recommended)

```bash
defenders conf
```

### Option 2: Environment variable

```bash
export ADO_PAT="your-pat-token"
```

### Option 3: Command-line flag

```bash
defenders release run <url> -t <token>
defenders pr --approve <url> -t <token>
```

**Priority order:** Flag > Environment variable > Config file

### Creating a PAT Token

1. Go to: https://dev.azure.com/{your-org}/_usersSettings/tokens
2. Create a new token with required scopes:
   - **Work Items**: Read & Write (for `cado`)
   - **Code**: Read & Write (for `prme`, `pr`)
   - **Build**: Read & Execute (for `release`)

---

## Examples

### Daily workflow

```bash
# Start of day: create a feature
defenders cado --title="JIRA-123: Implement login page" --parent=99999

# Work on code, then create PR
defenders prme -i 12345 -t "JIRA-123: Implement login page"

# Approve a colleague's PR
defenders pr --approve https://msazure.visualstudio.com/One/_git/MyRepo/pullrequest/456
```

### Release workflow

```bash
# Run PR validation pipeline, then trigger release when it succeeds
defenders release monitor-trigger \
  "https://msazure.visualstudio.com/One/_build/results?buildId=137041397" \
  "https://msazure.visualstudio.com/One/_build?definitionId=373994"
```

### Using another user's PAT

When a colleague needs you to approve their PR or run a release:

```bash
# They share their PAT with you
defenders pr --approve <pr-url> -t <their-pat>

# Or run a pipeline with their permissions
defenders release run <pipeline-url> -t <their-pat>
```

---

## Cross-Platform Support

This CLI works on **Windows**, **Linux**, and **macOS**.

### Build for different platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o defenders-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o defenders-mac

# Windows
GOOS=windows GOARCH=amd64 go build -o defenders.exe
```

### Requirements

- Azure CLI (`az`) installed and in PATH
- Git (for `prme` command)

---

## License

MIT