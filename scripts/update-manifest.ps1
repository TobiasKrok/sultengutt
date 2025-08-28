# Script to manually update Scoop manifest
# Usage: .\update-manifest.ps1 v1.0.0

param(
    [Parameter(Mandatory=$true)]
    [string]$Version
)

Write-Host "Fetching release information for $Version..." -ForegroundColor Green

# Remove 'v' prefix if present
$VersionClean = $Version -replace '^v', ''

# Get release data from GitHub API
try {
    $ReleaseUrl = "https://api.github.com/repos/TobiasKrok/sultengutt/releases/tags/$Version"
    $ReleaseData = Invoke-RestMethod -Uri $ReleaseUrl -Method Get
} catch {
    Write-Host "Error: Release $Version not found" -ForegroundColor Red
    exit 1
}

# Find Windows asset
$WindowsAsset = $ReleaseData.assets | Where-Object { $_.name -like "*windows-amd64.zip" }

if (-not $WindowsAsset) {
    Write-Host "Error: Could not find Windows binary in release" -ForegroundColor Red
    exit 1
}

$WindowsUrl = $WindowsAsset.browser_download_url

Write-Host "Downloading binary to calculate SHA256..." -ForegroundColor Yellow

# Download and calculate SHA256
$TempFile = Join-Path $env:TEMP "sultengutt-windows.zip"
Invoke-WebRequest -Uri $WindowsUrl -OutFile $TempFile

$Hash = Get-FileHash -Path $TempFile -Algorithm SHA256
$SHA256 = $Hash.Hash.ToLower()

# Clean up
Remove-Item $TempFile

# Generate manifest
$Manifest = @{
    version = $VersionClean
    description = "Cross-platform desktop reminder for ordering surprise dinners"
    homepage = "https://github.com/TobiasKrok/sultengutt"
    license = "MIT"
    architecture = @{
        "64bit" = @{
            url = $WindowsUrl
            hash = $SHA256
        }
    }
    bin = "sultengutt.exe"
    pre_uninstall = @(
        '& "$dir\sultengutt.exe" uninstall --confirm 2>$null | Out-Null'
    )
    checkver = @{
        github = "https://github.com/TobiasKrok/sultengutt"
    }
    autoupdate = @{
        architecture = @{
            "64bit" = @{
                url = 'https://github.com/TobiasKrok/sultengutt/releases/download/v$version/sultengutt-windows-amd64.zip'
            }
        }
    }
    post_install = @(
        "Write-Host 'To set up Sultengutt, run: sultengutt install' -ForegroundColor Green",
        "Write-Host 'To check status: sultengutt status' -ForegroundColor Green"
    )
}

# Convert to JSON and save
$JsonContent = $Manifest | ConvertTo-Json -Depth 10
$JsonContent | Out-File -FilePath "sultengutt.json" -Encoding UTF8

Write-Host "`nManifest generated: sultengutt.json" -ForegroundColor Green
Write-Host "`nDetails:" -ForegroundColor Cyan
Write-Host "  Version: $VersionClean"
Write-Host "  URL: $WindowsUrl"
Write-Host "  SHA256: $SHA256"
Write-Host "`nTo publish to scoop-bucket, copy this file to bucket/sultengutt.json in your bucket repository" -ForegroundColor Yellow