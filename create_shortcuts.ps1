$wshell = New-Object -ComObject WScript.Shell
$dir = (Get-Item -Path .).FullName
$binary = Join-Path $dir "stardf-anime.exe"
if (-not (Test-Path $binary)) {
    $binary = Join-Path $dir "build\stardf-anime.exe"
}
if (-not (Test-Path $binary)) {
    throw "Binário não encontrado. Gere o stardf-anime.exe antes de criar atalhos."
}

$icon = Join-Path $dir "mobile\flutter\windows\runner\resources\app_icon.ico"
if (-not (Test-Path $icon)) {
    $icon = $binary
}

# Console Shortcut
$shortcut = $wshell.CreateShortcut((Join-Path $dir "StarDF-Anime (Console).lnk"))
$shortcut.TargetPath = $binary
$shortcut.WorkingDirectory = $dir
$shortcut.Description = "StarDF-Anime - Console Mode"
$shortcut.IconLocation = $icon
$shortcut.Save()

# Web UI Shortcut
$shortcut2 = $wshell.CreateShortcut((Join-Path $dir "StarDF-Anime (Web UI).lnk"))
$shortcut2.TargetPath = $binary
$shortcut2.Arguments = "-web"
$shortcut2.WorkingDirectory = $dir
$shortcut2.Description = "StarDF-Anime - Web Mode"
$shortcut2.IconLocation = $icon
$shortcut2.Save()

Write-Host "Atalhos criados com sucesso na raiz do projeto!" -ForegroundColor Green
