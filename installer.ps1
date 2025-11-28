param(
  [string]$InstallRoot = "$env:LOCALAPPDATA\Programs\Clockwise",
  [string]$ModuleBase = 'codeberg.org/clockwise-lang/clockwise',
  [string]$RepoUrl
)

$ErrorActionPreference = 'Stop'

if (-not $RepoUrl) {
  $RepoUrl = "https://$ModuleBase.git"
}

function Info($m){ Write-Host "[INFO] $m" }
function Warn($m){ Write-Host "[WARN] $m" -ForegroundColor Yellow }
function Err($m){ Write-Host "[ERROR] $m" -ForegroundColor Red }

# 1) Check Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
  Err 'Go is not installed or not on PATH. Please install Go and re-run.'
  exit 1
}

# 2) Prepare dirs
$InstallBin = Join-Path $InstallRoot 'bin'
New-Item -ItemType Directory -Force -Path $InstallBin | Out-Null
$Tmp = Join-Path $env:TEMP ("cw_install_" + (Get-Random))
New-Item -ItemType Directory -Force -Path $Tmp | Out-Null

# Helper: copy binaries from $GOBIN/GOPATH/bin
function Copy-GoBins {
  $goBin = (go env GOBIN).Trim()
  if (-not $goBin) {
    $goPath = (go env GOPATH).Trim()
    if ($goPath) { $goBin = Join-Path $goPath 'bin' }
  }
  if (-not $goBin) { Err 'Unable to resolve Go bin directory.'; return $false }
  foreach ($b in 'cw.exe','cwfmt.exe','cwdoc.exe') {
    $src = Join-Path $goBin $b
    if (-not (Test-Path $src)) { Err "Expected $b not found in $goBin"; return $false }
    Copy-Item $src (Join-Path $InstallBin $b) -Force
  }
  return $true
}

# 3) Default path: go install
Info 'Installing Clockwise tools via `go install`...'
# Prefer direct fetching (skip proxy, helps in restricted networks)
$env:GOPROXY = 'direct'
$env:GOSUMDB = 'off'
$tools = 'cw','cwfmt','cwdoc'
$goInstallOk = $true
foreach ($p in $tools) {
  Info "go install $ModuleBase/cmd/$p@latest"
  & go install "$ModuleBase/cmd/$p@latest"
  if ($LASTEXITCODE -ne 0) { $goInstallOk = $false; break }
}
if ($goInstallOk) {
  if (-not (Copy-GoBins)) { $goInstallOk = $false }
  else {
    $cwPath = Join-Path $InstallBin 'cw.exe'
    $cwcPath = Join-Path $InstallBin 'cwc.exe'
    if (Test-Path $cwPath) {
      Copy-Item $cwPath $cwcPath -Force
      Info 'Created cwc.exe alias for cw.exe'
    }
  }
}

# 4) Fallback: build from local repo if available, else clone and patch
if (-not $goInstallOk) {
  if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Err 'Git is not available; cannot fall back to source build.'
    exit 1
  }
  $localRepo = $PSScriptRoot
  $hasLocal = Test-Path (Join-Path $localRepo 'cmd/cw/main.go')
  $repoDir = $null
  if ($hasLocal) {
    Info 'Attempting source build from local checkout (script directory)...'
    $repoDir = $localRepo
  } else {
    $repoDir = Join-Path $Tmp 'clockwise'
    Info "Cloning from $RepoUrl ..."
    & git clone --depth=1 $RepoUrl $repoDir | Out-Null
    if ($LASTEXITCODE -ne 0) { Err 'git clone failed'; exit 1 }
  }
  Push-Location $repoDir
  try {
    # Patch imports across the repo
    $legacy = 'github.com/your-org/clockwise'
    $repl = $ModuleBase
    $patchedCount = 0
    Get-ChildItem -Path $repoDir -Recurse -Include *.go | ForEach-Object {
      $path = $_.FullName
      $t = Get-Content $path -Raw
      $orig = $t
      # Replace legacy github module path with $ModuleBase
      if ($t -like "*${legacy}*") { $t = $t -replace [Regex]::Escape($legacy), $repl }
      # Replace non-existent ast package with parser package
      $t = $t -replace [Regex]::Escape("$ModuleBase/ast"), "$ModuleBase/parser"
      # Remove broken internal/logger import line if present
      $loggerPattern = '(?m)^\s*"' + [Regex]::Escape($ModuleBase) + '/internal/logger"\s*\r?\n'
      $t = $t -replace ($loggerPattern), ''
      if ($t -ne $orig) {
        $patchedCount++
        Set-Content -Path $path -NoNewline -Value $t
      }
    }
    if ($patchedCount -gt 0) { Info "Patched imports/paths in $patchedCount files" } else { Info 'No import patching required.' }

    # Try to tidy, but do not fail install if it can't reach network
    Info 'Running go mod tidy (best-effort)...'
    & go mod tidy
    if ($LASTEXITCODE -ne 0) { Warn 'go mod tidy failed; continuing with local build' }
    Info 'Building cw...'; & go build -o (Join-Path $InstallBin 'cw.exe') ./cmd/cw
    if ($LASTEXITCODE -ne 0) { throw 'go build cw failed' }
    $cwPath = Join-Path $InstallBin 'cw.exe'
    $cwcPath = Join-Path $InstallBin 'cwc.exe'
    if (Test-Path $cwPath) {
      Copy-Item $cwPath $cwcPath -Force
      Info 'Created cwc.exe alias for cw.exe'
    } else {
      # Attempt direct build of cwc.exe name from same package
      Info 'cw.exe not found; attempting to build cwc.exe directly...'
      & go build -o $cwcPath ./cmd/cw
      if ($LASTEXITCODE -ne 0) { throw 'go build cwc failed' }
    }
    Info 'Building cwfmt...'; & go build -o (Join-Path $InstallBin 'cwfmt.exe') ./cmd/cwfmt
    if ($LASTEXITCODE -ne 0) { throw 'go build cwfmt failed' }
    Info 'Building cwdoc...'; & go build -o (Join-Path $InstallBin 'cwdoc.exe') ./cmd/cwdoc
    if ($LASTEXITCODE -ne 0) { throw 'go build cwdoc failed' }
  } catch {
    Err "Source build failed: $($_.Exception.Message)"; Pop-Location; exit 1
  }
  Pop-Location
}

# 5) Update user PATH
try {
  $p = [Environment]::GetEnvironmentVariable('Path','User')
  if (-not $p) { $p = '' }
  $ib = [IO.Path]::GetFullPath($InstallBin)
  if ($p -notlike ('*' + $ib + '*')) {
    if ($p -and -not $p.EndsWith(';')) { $p += ';' }
    $p += $ib
    [Environment]::SetEnvironmentVariable('Path',$p,'User')
    Info 'PATH updated.'
  } else {
    Info 'PATH already contains install dir.'
  }
} catch {
  Warn 'Failed to update user PATH automatically. Please add it manually.'
}

# 6) Validate and Done
$cwOk = Test-Path (Join-Path $InstallBin 'cw.exe')
$cwcOk = Test-Path (Join-Path $InstallBin 'cwc.exe')
if (-not ($cwOk -or $cwcOk)) {
  Err 'Main CLI (cw.exe/cwc.exe) was not built. Installation aborted.'
  exit 1
}
Info "Clockwise installed to: $InstallBin"
Info 'Binaries: cw.exe, cwc.exe, cwfmt.exe, cwdoc.exe'
Write-Host '[TIP] Open a new terminal and run: cwc --help'

Remove-Item -Recurse -Force $Tmp -ErrorAction SilentlyContinue
