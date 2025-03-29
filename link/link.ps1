param (
    [Parameter(Position=0)]
    [ValidateSet("status", "list", "create", "remove")]
    [string]$Mode = "list",
    
    [Parameter(Position=1)]
    [string]$SoftwareFilter = ""
)

# Global Configuration 
$global:Config = @{
    CleanupScript = "place_holder_clean_up.ps1"  # Cleanup script to run after removal
    User = $env:USERNAME
    ScriptDriveLetter = (Get-Location).Drive.Name + ":"
    PathsJsonFile = "$ScriptDriveLetter\system\config\junctions.json"  
}

$ErrorActionPreference = 'Stop'

function Check-SoftwareExists {
    param (
        [string]$SoftwareName,
        [object]$PathMappings
    )
    
    if ([string]::IsNullOrEmpty($SoftwareName)) {
        return $true  #  processing all software
    }
    
    return $PathMappings.PSObject.Properties.Name -contains $SoftwareName
}

function Show-LinkDetailedStatus {
    param ([string]$SoftwareFilter = $SoftwareFilter)

    $pathMappings = Get-Content -Raw -Path $global:Config.PathsJsonFile | ConvertFrom-Json
    
    if ($SoftwareFilter -and -not (Check-SoftwareExists -SoftwareName $SoftwareFilter -PathMappings $pathMappings)) {
        Write-Host "Error: Software '$SoftwareFilter' doesn't exist in configuration file" -ForegroundColor Red
        exit 1
    }

    foreach ($software in $pathMappings.PSObject.Properties) {
        if ($SoftwareFilter -and $software.Name -ne $SoftwareFilter) { continue }

        $softwareName = $software.Name
        Write-Host $softwareName -ForegroundColor Cyan
        $missingTargets = @()

        foreach ($path in $software.Value.Paths) {
            $sourcePath = $path.SourcePath -replace '\$user', $global:Config.User
            $targetPath = $path.TargetPath -replace '\$scriptDriveLetter', $global:Config.ScriptDriveLetter

            if (-not (Test-Path $sourcePath)) {
                Write-Host " - missing " -NoNewline -ForegroundColor Red
                Write-Host ""$sourcePath"" -ForegroundColor Red
                continue
            }

            $linkItem = Get-Item $sourcePath
            if ($linkItem.Attributes -match 'ReparsePoint') {
                $currentTarget = $linkItem.Target
                if ($currentTarget -eq $targetPath) {
                    Write-Host " - ok " -NoNewline -ForegroundColor Green
                    Write-Host ""$sourcePath" ===> "$targetPath"" -ForegroundColor White

                    if (-not (Test-Path $targetPath)) {
                        $missingTargets += $targetPath
                    }
                }
                else {
                    Write-Host " - missconf " -NoNewline -ForegroundColor Yellow
                    Write-Host ""$sourcePath" ===> "$currentTarget"" -ForegroundColor Red
                }
            }
            else {
                Write-Host " - dir " -NoNewline -ForegroundColor Blue
                Write-Host ""$sourcePath"" -ForegroundColor White
            }
        }

        if ($missingTargets.Count -gt 0) {
            Write-Host "  [warning]: Target directory does not exist:" -ForegroundColor Yellow
            foreach ($target in $missingTargets) {
                Write-Host "      Path: $target" -ForegroundColor Yellow
            }
            Write-Host "`n    Create directory or run script with 'create' mode to fix" -ForegroundColor Yellow
        }
        Write-Host ""
    }
}

function Show-LinkSummary {
    $pathMappings = Get-Content -Raw -Path $global:Config.PathsJsonFile | ConvertFrom-Json
    
    if ($SoftwareFilter -and -not (Check-SoftwareExists -SoftwareName $SoftwareFilter -PathMappings $pathMappings)) {
        Write-Host "Error: Software '$SoftwareFilter' doesn't exist in configuration file" -ForegroundColor Red
        exit 1
    }

    $maxNameLength = ($pathMappings.PSObject.Properties.Name | Measure-Object -Maximum -Property Length).Maximum

    foreach ($software in $pathMappings.PSObject.Properties) {
        if ($SoftwareFilter -and $software.Name -ne $SoftwareFilter) { continue }
        
        $softwareName = $software.Name
        $totalPaths = $software.Value.Paths.Count
        $validLinks = 0

        foreach ($path in $software.Value.Paths) {
            $sourcePath = $path.SourcePath -replace '\$user', $global:Config.User
            $targetPath = $path.TargetPath -replace '\$scriptDriveLetter', $global:Config.ScriptDriveLetter

            if (Test-Path $sourcePath) {
                $linkItem = Get-Item $sourcePath
                if ($linkItem.Attributes -match 'ReparsePoint') {
                    $currentTarget = $linkItem.Target
                    if ($currentTarget -eq $targetPath -and (Test-Path $targetPath)) {
                        $validLinks++
                    }
                }
            }
        }

        $status = if ($validLinks -eq $totalPaths) { "ok" } else { "missing" }
        $dots = "." * ($maxNameLength - $softwareName.Length + 3)

        $color = if ($status -eq "ok") { "Green" } else { "Red" }
        Write-Host "$softwareName$dots" -NoNewline
        Write-Host ("[{0}/{1}]" -f $validLinks, $totalPaths) -NoNewline -ForegroundColor $color
        Write-Host "[$status]" -ForegroundColor $color
    }
}

function Create-OrUpdateLinks {
    param ([string]$SoftwareFilter = $SoftwareFilter)

    $pathMappings = Get-Content -Raw -Path $global:Config.PathsJsonFile | ConvertFrom-Json
    
    if ($SoftwareFilter -and -not (Check-SoftwareExists -SoftwareName $SoftwareFilter -PathMappings $pathMappings)) {
        Write-Host "Error: Software '$SoftwareFilter' doesn't exist in configuration file" -ForegroundColor Red
        exit 1
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"

    foreach ($software in $pathMappings.PSObject.Properties) {
        if ($SoftwareFilter -and $software.Name -ne $SoftwareFilter) { continue }

        $softwareName = $software.Name
        Write-Host "Processing $softwareName" -ForegroundColor Cyan

        foreach ($path in $software.Value.Paths) {
            $sourcePath = $path.SourcePath -replace '\$user', $global:Config.User
            $targetPath = $path.TargetPath -replace '\$scriptDriveLetter', $global:Config.ScriptDriveLetter

            if (-not (Test-Path $targetPath)) {
                Write-Host "Creating target directory: $targetPath" -ForegroundColor DarkGray
                New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
            }

            if (-not (Test-Path $sourcePath)) {
                Write-Host "Creating link: "$sourcePath" -> "$targetPath"" -ForegroundColor Green
                New-Item -ItemType Junction -Path $sourcePath -Target $targetPath | Out-Null
            }
            else {
                $linkItem = Get-Item $sourcePath
                if ($linkItem.Attributes -match 'ReparsePoint') {
                    $currentTarget = $linkItem.Target
                    if ($currentTarget -eq $targetPath) {
                        Write-Host "Link ok: "$sourcePath" already points to "$targetPath"" -ForegroundColor Green
                    }
                    else {
                        Write-Host "Recreating misconfigured link: "$sourcePath" -> "$targetPath"" -ForegroundColor Yellow
                        Remove-Item $sourcePath -Force
                        New-Item -ItemType Junction -Path $sourcePath -Target $targetPath | Out-Null
                    }
                }
                else {
                    $subFolderName = Split-Path $sourcePath -Leaf
                    $backupRoot = "$($global:Config.ScriptDriveLetter)\Apps\backup\$softwareName\data-$timestamp"
                    $backupDir = Join-Path $backupRoot $subFolderName
                    New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
                    Get-ChildItem -Path $sourcePath -Force | Move-Item -Destination $backupDir -Force
                    Remove-Item $sourcePath -Force
                    New-Item -ItemType Junction -Path $sourcePath -Target $targetPath | Out-Null
                    Write-Host "Replaced directory with link: "$sourcePath" -> "$targetPath"" -ForegroundColor Green
                }
            }
        }
        Write-Host ""
    }
}

function Remove-Links {
    param ([string]$SoftwareFilter = $SoftwareFilter)

    $pathMappings = Get-Content -Raw -Path $global:Config.PathsJsonFile | ConvertFrom-Json
    
    if ($SoftwareFilter -and -not (Check-SoftwareExists -SoftwareName $SoftwareFilter -PathMappings $pathMappings)) {
        Write-Host "Error: Software '$SoftwareFilter' doesn't exist in configuration file" -ForegroundColor Red
        exit 1
    }

    foreach ($software in $pathMappings.PSObject.Properties) {
        if ($SoftwareFilter -and $software.Name -ne $SoftwareFilter) { continue }

        $softwareName = $software.Name
        Write-Host "Removing links for: $softwareName" -ForegroundColor Cyan

        foreach ($path in $software.Value.Paths) {
            $sourcePath = $path.SourcePath -replace '\$user', $global:Config.User
            $targetPath = $path.TargetPath -replace '\$scriptDriveLetter', $global:Config.ScriptDriveLetter

            if (Test-Path $sourcePath) {
                $item = Get-Item $sourcePath -Force
                if ($item.Attributes -match 'ReparsePoint') {
                    Write-Host " - Removing junction: " -NoNewline -ForegroundColor Green
                    Write-Host ""$sourcePath"" -ForegroundColor White
                    Remove-Item $sourcePath -Force
                }
                else {
                    Write-Host " - Removing directory: " -NoNewline -ForegroundColor Yellow
                    Write-Host ""$sourcePath"" -ForegroundColor White
                    Remove-Item $sourcePath -Recurse -Force
                }
            }
            else {
                Write-Host " - Not found: " -NoNewline -ForegroundColor DarkGray
                Write-Host ""$sourcePath"" -ForegroundColor White
            }
        }
        Write-Host ""
    }

    # Cleanup script execution (uncomment to enable)
    <#
    if (Test-Path $global:Config.CleanupScript) {
        Start-Process powershell.exe -ArgumentList "-File `"$($global:Config.CleanupScript)`"" -NoNewWindow
    }
    #>
}

switch ($Mode.ToLower()) {
    "status" { Show-LinkDetailedStatus }
    "list"   { Show-LinkSummary }
    "create" { Create-OrUpdateLinks }
    "remove" { Remove-Links }
    default  { 
        Write-Host "Specify either 'status', 'list', 'create', or 'remove'" -ForegroundColor Yellow 
        exit 1
    }
}