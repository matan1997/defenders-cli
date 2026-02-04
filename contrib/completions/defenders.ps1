# PowerShell completion for defenders CLI

Register-ArgumentCompleter -Native -CommandName defenders -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @{
        'conf' = 'Configuration management'
        'get-token' = 'Get authentication token'
        'cado' = 'Azure DevOps operations'
        'prme' = 'Create pull request'
        'release' = 'Run pipeline/release'
        'pr' = 'Pull request handler'
        'help' = 'Show help information'
    }

    $confOptions = @{
        '-s' = 'Show current configuration'
        '--show' = 'Show current configuration'
        '-p' = 'Show configuration file path'
        '--path' = 'Show configuration file path'
        '-r' = 'Reset configuration'
        '--reset' = 'Reset configuration'
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    $cadoOptions = @{
        '-u' = 'Azure DevOps URL'
        '--url' = 'Azure DevOps URL'
        '-j' = 'JSON data'
        '--json' = 'JSON data'
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    $prmeOptions = @{
        '-t' = 'PR title'
        '--title' = 'PR title'
        '-d' = 'PR description'
        '--description' = 'PR description'
        '-s' = 'Source branch'
        '--source' = 'Source branch'
        '-b' = 'Target branch'
        '--target' = 'Target branch'
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    $releaseOptions = @{
        '-u' = 'Pipeline URL'
        '--url' = 'Pipeline URL'
        '-m' = 'Monitor and trigger'
        '--monitor' = 'Monitor and trigger'
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    $prOptions = @{
        '-u' = 'Pull request URL'
        '--url' = 'Pull request URL'
        '-a' = 'Approve PR'
        '--approve' = 'Approve PR'
        '-c' = 'Complete PR'
        '--complete' = 'Complete PR'
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    $helpFlags = @{
        '-h' = 'Show help'
        '--help' = 'Show help'
    }

    # Get all words in the command line
    $words = $commandAst.ToString() -split '\s+' | Where-Object { $_ }

    # If we're at the first position (just after 'defenders')
    if ($words.Count -le 2) {
        $commands.GetEnumerator() | Where-Object { $_.Key -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new(
                $_.Key,
                $_.Key,
                'ParameterValue',
                $_.Value
            )
        }
        return
    }

    # Get the subcommand
    $subcommand = $words[1]

    # Complete options based on subcommand
    $options = switch ($subcommand) {
        'conf' { $confOptions }
        'get-token' { $helpFlags }
        'cado' { $cadoOptions }
        'prme' { $prmeOptions }
        'release' { $releaseOptions }
        'pr' { $prOptions }
        'help' { 
            $commands.GetEnumerator() | Where-Object { $_.Key -like "$wordToComplete*" -and $_.Key -ne 'help' } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new(
                    $_.Key,
                    $_.Key,
                    'ParameterValue',
                    $_.Value
                )
            }
            return
        }
        default { @{} }
    }

    # Return matching options
    $options.GetEnumerator() | Where-Object { $_.Key -like "$wordToComplete*" } | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new(
            $_.Key,
            $_.Key,
            'ParameterName',
            $_.Value
        )
    }
}
