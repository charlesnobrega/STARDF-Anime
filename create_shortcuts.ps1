$wshell = New-Object -ComObject WScript.Shell
$dir = (Get-Item -Path .).FullName
$binary = Join-Path $dir "build\stardf-anime.exe"

# Console Shortcut
$shortcut = $wshell.CreateShortcut((Join-Path $dir "StarDF-Anime (Console).lnk"))
$shortcut.TargetPath = $binary
$shortcut.WorkingDirectory = $dir
$shortcut.Description = "StarDF-Anime - Console Mode"
$shortcut.Save()

# Web UI Shortcut
$shortcut2 = $wshell.CreateShortcut((Join-Path $dir "StarDF-Anime (Web UI).lnk"))
$shortcut2.TargetPath = $binary
$shortcut2.Arguments = "-web"
$shortcut2.WorkingDirectory = $dir
$shortcut2.Description = "StarDF-Anime - Web Mode"
$shortcut2.Save()

Write-Host "Atalhos criados com sucesso na raiz do projeto!" -ForegroundColor Green
