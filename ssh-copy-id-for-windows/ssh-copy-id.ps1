param(
    [string]$Identity,
    [string]$User,
    [string]$Hostname,
    [int]$Port = 22,
    [switch]$h,
    [switch]$help
)

function Show-Usage {
    Write-Host "SSH Copy ID for Windows - Copy local SSH public key to remote server"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  .\ssh-copy-id.ps1 [-Identity keypath] [-User username] [-Hostname host] [-Port port]"
    Write-Host ""
    Write-Host "Parameters:"
    Write-Host "  -Identity  SSH key file path"
    Write-Host "  -User      Remote server username"
    Write-Host "  -Hostname  Remote server IP or hostname"
    Write-Host "  -Port      SSH port (default: 22)"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\ssh-copy-id.ps1 -User ubuntu -Hostname 192.168.1.100"
}

function Test-PowerShellVersion {
    $version = $PSVersionTable.PSVersion
    $majorVersion = $version.Major

    Write-Host "PowerShell Version Check:" -ForegroundColor Cyan
    Write-Host "  Current Version: PowerShell $($version.Major).$($version.Minor)" -ForegroundColor Green

    if ($majorVersion -eq 5) {
        Write-Host "  Type: Windows PowerShell 5.x" -ForegroundColor Yellow
        Write-Host "  Suggestion: Consider upgrading to PowerShell 7+ for better performance" -ForegroundColor Yellow
    } elseif ($majorVersion -ge 7) {
        Write-Host "  Type: PowerShell Core 7+" -ForegroundColor Green
        Write-Host "  Status: Modern version with excellent performance" -ForegroundColor Green
    } else {
        Write-Host "  Warning: Version too old, may have compatibility issues" -ForegroundColor Red
    }

    Write-Host ""
    return $majorVersion
}

function Test-SshCommand {
    try {
        $null = Get-Command ssh -ErrorAction Stop
        return $true
    }
    catch {
        return $false
    }
}

function Get-DefaultKeyPath {
    $sshDir = Join-Path $env:USERPROFILE ".ssh"
    $keyTypes = @("id_rsa.pub", "id_ed25519.pub", "id_ecdsa.pub", "id_dsa.pub")

    foreach ($keyType in $keyTypes) {
        $keyPath = Join-Path $sshDir $keyType
        if (Test-Path $keyPath) {
            return $keyPath
        }
    }

    return $null
}

function Read-SecurePassword {
    param([string]$Prompt)

    Write-Host $Prompt -NoNewline
    $password = ""
    do {
        $key = [System.Console]::ReadKey($true)
        if ($key.Key -eq "Enter") {
            Write-Host ""
            break
        }
        elseif ($key.Key -eq "Backspace") {
            if ($password.Length -gt 0) {
                $password = $password.Substring(0, $password.Length - 1)
                Write-Host "`b `b" -NoNewline
            }
        }
        else {
            $password += $key.KeyChar
            Write-Host "*" -NoNewline
        }
    } while ($true)

    return $password
}

function Copy-SSHKey {
    param(
        [string]$KeyPath,
        [string]$Username,
        [string]$Hostname,
        [int]$SshPort,
        [string]$Password
    )

    try {
        Write-Host "Reading public key file: $KeyPath" -ForegroundColor Yellow

        if (-not (Test-Path $KeyPath)) {
            throw "Public key file does not exist: $KeyPath"
        }

        $publicKey = Get-Content $KeyPath -Raw
        $publicKey = $publicKey.Trim()

        if ([string]::IsNullOrEmpty($publicKey)) {
            throw "Public key file is empty or cannot be read"
        }

        Write-Host "Public key content: $($publicKey.Substring(0, [Math]::Min(50, $publicKey.Length)))..." -ForegroundColor Green

        # Create remote command
        $cmd1 = "mkdir -p ~/.ssh"
        $cmd2 = "chmod 700 ~/.ssh"
        $cmd3 = "echo '$publicKey' >> ~/.ssh/authorized_keys"
        $cmd4 = "chmod 600 ~/.ssh/authorized_keys"
        $cmd5 = "echo 'SSH key added successfully'"
        $remoteCommand = "$cmd1 && $cmd2 && $cmd3 && $cmd4 && $cmd5"

        Write-Host "Connecting to remote server: $Username@$Hostname`:$SshPort" -ForegroundColor Yellow

        # Check if plink is available
        $plinkAvailable = $false
        try {
            $null = Get-Command plink -ErrorAction Stop
            $plinkAvailable = $true
        }
        catch {
            $plinkAvailable = $false
        }

        if ($plinkAvailable) {
            Write-Host "Using PuTTY plink for SSH connection..." -ForegroundColor Yellow

            $plinkArgs = @(
                "-ssh",
                "-P", $SshPort,
                "-l", $Username,
                "-pw", $Password,
                "-batch",
                $Hostname,
                $remoteCommand
            )

            try {
                $output = & plink @plinkArgs 2>&1
                $exitCode = $LASTEXITCODE

                Write-Host $output -ForegroundColor Green

                if ($exitCode -eq 0) {
                    Write-Host "`nSuccess: SSH key copied successfully!" -ForegroundColor Green
                    Write-Host "You can now log in without password using:" -ForegroundColor Cyan
                    Write-Host "ssh -p $SshPort $Username@$Hostname" -ForegroundColor Cyan
                    return $true
                }
                else {
                    Write-Host "`nFailed: SSH key copy failed (Exit code: $exitCode)" -ForegroundColor Red
                    return $false
                }
            }
            catch {
                Write-Host "`nError: plink error: $($_.Exception.Message)" -ForegroundColor Red
                return $false
            }
        }
        else {
            Write-Host "Using simple SSH connection method..." -ForegroundColor Yellow
            Write-Host "Note: May require manual password entry" -ForegroundColor Yellow

            try {
                Write-Host "Executing SSH command. You may need to enter password manually if prompted..." -ForegroundColor Cyan

                $sshArgs = @(
                    "-o", "StrictHostKeyChecking=no",
                    "-o", "UserKnownHostsFile=NUL",
                    "-p", $SshPort,
                    "$Username@$Hostname",
                    $remoteCommand
                )

                $process = Start-Process -FilePath "ssh" -ArgumentList $sshArgs -Wait -PassThru -NoNewWindow

                if ($process.ExitCode -eq 0) {
                    Write-Host "`nSuccess: SSH key copied successfully!" -ForegroundColor Green
                    Write-Host "You can now log in without password using:" -ForegroundColor Cyan
                    Write-Host "ssh -p $SshPort $Username@$Hostname" -ForegroundColor Cyan
                    return $true
                }
                else {
                    Write-Host "`nFailed: SSH key copy failed (Exit code: $($process.ExitCode))" -ForegroundColor Red
                    Write-Host "Try installing PuTTY (plink) for better password automation" -ForegroundColor Yellow
                    return $false
                }
            }
            catch {
                Write-Host "`nError: SSH execution failed: $($_.Exception.Message)" -ForegroundColor Red
                Write-Host "Suggestions:" -ForegroundColor Yellow
                Write-Host "1. Install PuTTY (includes plink command)" -ForegroundColor Yellow
                Write-Host "2. Or manually use ssh-copy-id if available" -ForegroundColor Yellow
                return $false
            }
        }
    }
    catch {
        Write-Host "`nError: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Main program starts here
if ($h -or $help) {
    Show-Usage
    exit 0
}

Write-Host "=== SSH Copy ID for Windows ===" -ForegroundColor Cyan
Write-Host ""

$psVersion = Test-PowerShellVersion

if (-not (Test-SshCommand)) {
    Write-Host "Error: SSH command not found. Please ensure OpenSSH client is installed." -ForegroundColor Red
    Write-Host "You can install it via:" -ForegroundColor Yellow
    Write-Host "1. Windows 10/11: Settings > Apps > Optional Features > Add Feature > OpenSSH Client" -ForegroundColor Yellow
    Write-Host "2. Or install Git for Windows (includes SSH)" -ForegroundColor Yellow
    exit 1
}

if ([string]::IsNullOrEmpty($Identity)) {
    $defaultKey = Get-DefaultKeyPath
    if ($defaultKey) {
        $Identity = $defaultKey
        Write-Host "Using default key file: $Identity" -ForegroundColor Green
    }
    else {
        Write-Host "No default SSH key file found." -ForegroundColor Yellow
        $Identity = Read-Host "Please enter SSH public key file path"
    }
}

if ([string]::IsNullOrEmpty($User)) {
    $User = Read-Host "Please enter remote server username"
}

if ([string]::IsNullOrEmpty($Hostname)) {
    $Hostname = Read-Host "Please enter remote server IP address or hostname"
}

if ($Port -eq 22) {
    $portInput = Read-Host "Please enter SSH port (default: 22)"
    if (-not [string]::IsNullOrEmpty($portInput)) {
        $Port = [int]$portInput
    }
}

$password = Read-SecurePassword "Please enter password for $User@${Hostname}: "

if ([string]::IsNullOrEmpty($password)) {
    Write-Host "`nError: Password cannot be empty" -ForegroundColor Red
    exit 1
}

Write-Host "`nStarting SSH key copy..." -ForegroundColor Yellow
$success = Copy-SSHKey -KeyPath $Identity -Username $User -Hostname $Hostname -SshPort $Port -Password $password

if ($success) {
    Write-Host "`nOperation completed!" -ForegroundColor Green
}
else {
    Write-Host "`nOperation failed!" -ForegroundColor Red
    exit 1
}