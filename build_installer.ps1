<#
Build script to only build the installer binary. This satisfies the project policy
that the installer is the primary artifact and it will bootstrap/build other tools.
#>
param(
    [string]$Out = "bin\installer.exe"
)

Write-Host "Building installer..."
New-Item -ItemType Directory -Force -Path bin | Out-Null
go build -o $Out ./cmd/installer
if ($LASTEXITCODE -ne 0) {
    Write-Error "go build failed"
    exit $LASTEXITCODE
}
Write-Host "Built:" $Out
