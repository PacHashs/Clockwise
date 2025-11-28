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
$tools = 'cw','cwfmt','cwdoc'
$goInstallOk = $true
foreach ($p in $tools) {
  Info "go install $ModuleBase/cmd/$p@latest"
  try { & go install "$ModuleBase/cmd/$p@latest" } catch { $goInstallOk = $false; break }
}
if ($goInstallOk) {
  if (-not (Copy-GoBins)) { $goInstallOk = $false }
}

# 4) Fallback: clone and build from source, patching legacy imports if present
if (-not $goInstallOk) {
  if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Err 'Git is not available; cannot fall back to source build.'
    exit 1
  }
  $repoDir = Join-Path $Tmp 'clockwise'
  Info "Cloning from $RepoUrl ..."
  & git clone --depth=1 $RepoUrl $repoDir | Out-Null
  Push-Location $repoDir
  try {
    # Patch legacy imports across the whole repo (zince imz lazy to do so)
    $legacy = 'github.com/your-org/clockwise'
    $repl = $ModuleBase
    $patchedCount = 0
    Get-ChildItem -Path $repoDir -Recurse -Include *.go | ForEach-Object {
      $t = Get-Content $_.FullName -Raw
      if ($t -like "*${legacy}*") {
        ($t -replace [Regex]::Escape($legacy), $repl) | Set-Content -Path $_.FullName -NoNewline
        $patchedCount++
      }
    }
    if ($patchedCount -gt 0) { Info "Patched imports in $patchedCount files" } else { Info 'No import patching required.' }

    Info 'Running go mod tidy...'
    & go mod tidy
    Info 'Building cw...'; & go build -o (Join-Path $InstallBin 'cw.exe') ./cmd/cw
    Info 'Building cwfmt...'; & go build -o (Join-Path $InstallBin 'cwfmt.exe') ./cmd/cwfmt
    Info 'Building cwdoc...'; & go build -o (Join-Path $InstallBin 'cwdoc.exe') ./cmd/cwdoc
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

# 6) Done
Info "Clockwise installed to: $InstallBin"
Info 'Binaries: cw.exe, cwfmt.exe, cwdoc.exe'
Write-Host '[TIP] Open a new terminal and run: cw --help'

Remove-Item -Recurse -Force $Tmp -ErrorAction SilentlyContinue
