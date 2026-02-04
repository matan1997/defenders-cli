# Fish completion for defenders CLI
# Install: sudo cp defenders.fish /usr/share/fish/vendor_completions.d/

# Disable file completion for defenders
complete -c defenders -f

# Main subcommands - show when first argument position
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a conf -d "Configuration management"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a get-token -d "Get authentication token"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a cado -d "Create ADO Feature work item"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a prme -d "Create PR from current branch"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a release -d "Pipeline runner"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a pr -d "Pull request actions"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -a help -d "Show help information"
complete -c defenders -n "test (count (commandline -opc)) -eq 1" -s h -l help -d "Show help information"

# conf subcommand - has subcommands: show, path, reset, setup
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'conf'" -a show -d "Show current configuration"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'conf'" -a path -d "Show configuration file path"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'conf'" -a reset -d "Reset configuration to defaults"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'conf'" -s h -l help -d "Show help"

# get-token subcommand options
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'get-token'" -s h -l help -d "Show help"

# cado subcommand options
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'cado'" -l title -d "Title of the Feature work item" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'cado'" -l parent -d "Parent work item ID to link" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'cado'" -l assigned-to -d "Override assigned-to from config" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'cado'" -s h -l help -d "Show help"

# prme subcommand options
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'prme'" -s i -l work-item -d "Work item ID to link to the PR" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'prme'" -s t -l title -d "Custom PR title" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'prme'" -s h -l help -d "Show help"

# release subcommand - has subcommands: run, monitor-trigger
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'release'" -a run -d "Run a pipeline directly"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'release'" -a monitor-trigger -d "Monitor pipeline and trigger another"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'release'" -s t -l token -d "Personal Access Token" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'release'" -s i -l interval -d "Check interval in seconds" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'release'" -s h -l help -d "Show help"

# pr subcommand options
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'pr'" -l approve -d "Approve the Pull Request"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'pr'" -l reset -d "Reset your vote on the PR"
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'pr'" -s t -l token -d "Personal Access Token" -r
complete -c defenders -n "test (count (commandline -opc)) -ge 2; and test (commandline -opc)[2] = 'pr'" -s h -l help -d "Show help"
