 # Shell Completions for Defenders CLI

This directory contains shell completion scripts for the `defenders` CLI tool across multiple shells and platforms.

## Available Completions

- **Bash** (`defenders.bash`) - Linux/macOS/Windows Git Bash
- **Zsh** (`_defenders`) - macOS default shell (since Catalina)
- **Fish** (`defenders.fish`) - Fish shell
- **PowerShell** (`defenders.ps1`) - Windows PowerShell/PowerShell Core

## Installation Instructions

### Bash

#### Linux/macOS (User-specific)
```bash
mkdir -p ~/.local/share/bash-completion/completions
cp defenders.bash ~/.local/share/bash-completion/completions/defenders
```

#### Linux (System-wide - requires sudo)
```bash
sudo cp defenders.bash /etc/bash_completion.d/defenders
```

#### macOS with Homebrew (System-wide)
```bash
cp defenders.bash $(brew --prefix)/etc/bash_completion.d/defenders
```

After installation, restart your shell or run:
```bash
source ~/.bashrc  # or ~/.bash_profile on macOS
```

### Zsh

#### User-specific
```bash
mkdir -p ~/.zsh/completion
cp _defenders ~/.zsh/completion/
```

Add to your `~/.zshrc`:
```zsh
fpath=(~/.zsh/completion $fpath)
autoload -Uz compinit && compinit
```

#### System-wide (requires sudo)
```bash
sudo cp _defenders /usr/local/share/zsh/site-functions/
```

After installation, restart your shell or run:
```bash
exec zsh
```

### Fish

#### User-specific (Recommended)
```bash
mkdir -p ~/.config/fish/completions
cp defenders.fish ~/.config/fish/completions/
```

#### System-wide (requires sudo)
```bash
sudo cp defenders.fish /usr/share/fish/vendor_completions.d/
```

Completions work immediately - no need to restart Fish!

### PowerShell

#### Windows PowerShell / PowerShell Core

Add to your PowerShell profile:
```powershell
# Find your profile location
$PROFILE

# Edit your profile (creates if doesn't exist)
notepad $PROFILE

# Add this line to the profile
. "C:\path\to\defenders-cli\contrib\completions\defenders.ps1"
```

Alternatively, for user-specific installation:
```powershell
# Create directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "$HOME\.config\powershell\scripts"

# Copy completion script
Copy-Item defenders.ps1 "$HOME\.config\powershell\scripts\"

# Add to profile
Add-Content $PROFILE ". `"$HOME\.config\powershell\scripts\defenders.ps1`""
```

After installation, restart PowerShell or run:
```powershell
. $PROFILE
```

## Testing Completions

After installation, test the completions by typing:

```bash
defenders <TAB>
```

You should see command suggestions like:
- `conf` - Configuration management
- `get-token` - Get authentication token
- `cado` - Azure DevOps operations
- `prme` - Create pull request
- `release` - Run pipeline/release
- `pr` - Pull request handler

## Features

All completion scripts provide:
- ✅ Command completion
- ✅ Subcommand completion
- ✅ Flag/option completion
- ✅ Help text for commands (where supported by shell)
- ✅ Context-aware completions

## Troubleshooting

### Bash
If completions don't work, ensure bash-completion is installed:
```bash
# Ubuntu/Debian
sudo apt install bash-completion

# macOS
brew install bash-completion@2
```

### Zsh
If completions don't work, verify `compinit` is loaded in your `~/.zshrc`:
```zsh
autoload -Uz compinit && compinit
```

### Fish
Fish loads completions automatically from standard directories. No additional configuration needed.

### PowerShell
If completions don't work, ensure script execution is enabled:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Contributing

If you find issues with the completions or want to add support for additional flags, please submit an issue or pull request to the main repository.
