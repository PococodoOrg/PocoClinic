# Fail if git would stage secrets, PHI, or local env files. Run before commit.
$ErrorActionPreference = 'Stop'
Set-Location (Join-Path $PSScriptRoot '..')

$patterns = @(
    '\.env$',
    '[/\\]setup\.env$',
    '[/\\]personas\.env$',
    '\.db$',
    '\.db-wal$',
    '\.db-shm$',
    'bootstrap-admin-once\.txt',
    '[/\\]backend[/\\]data[/\\]',
    '^data[/\\]documents[/\\]',
    'pococlinic-backup-.*\.tar\.gz'
)

$candidates = @(
    git diff --cached --name-only --diff-filter=ACMRT
    git ls-files --others --exclude-standard
) | Sort-Object -Unique

if (-not $candidates) {
    Write-Host 'check-no-secrets: nothing to scan.'
    exit 0
}

$failed = $false
foreach ($path in $candidates) {
    foreach ($pattern in $patterns) {
        if ($path -match $pattern) {
            Write-Host "BLOCKED: $path (matches /$pattern/)"
            $failed = $true
        }
    }
    if (Test-Path $path -PathType Leaf) {
        $content = Get-Content $path -Raw -ErrorAction SilentlyContinue
        if ($content -match 'BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY') {
            Write-Host "BLOCKED: $path (contains a private key block)"
            $failed = $true
        }
    }
}

if ($failed) {
    Write-Host ''
    Write-Host 'check-no-secrets: remove blocked files from the commit or update .gitignore.'
    exit 1
}

Write-Host "check-no-secrets: OK ($($candidates.Count) paths scanned)"
