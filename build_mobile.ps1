# StarDF-Anime Mobile - Local Build Script
Write-Host "--- StarDF-Anime Mobile Build Automator ---" -ForegroundColor Cyan

# 1. Check prerequisites
$prereqs = @("go", "flutter", "gomobile")
foreach ($cmd in $prereqs) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Error "Prerequisite not found: $cmd"
        if ($cmd -eq "gomobile") {
            Write-Host "Install it with: go install golang.org/x/mobile/cmd/gomobile@latest" -ForegroundColor Yellow
        }
        exit 1
    }
}

# 2. Build Go Core Binding
Write-Host "`n[1/3] Building Go Core Binding (.aar)..." -ForegroundColor Green
Set-Location -Path "mobile"
& gomobile bind -v -target=android -o ./android/stardf_core.aar ./bridge.go
if ($LASTEXITCODE -ne 0) {
    Write-Error "Gomobile bind failed."
    exit 1
}

# 3. Flutter Pub Get
Write-Host "`n[2/3] Fetching Flutter dependencies..." -ForegroundColor Green
Set-Location -Path "flutter"
& flutter pub get
if ($LASTEXITCODE -ne 0) {
    Write-Error "Flutter pub get failed."
    exit 1
}

# 4. Build APK
Write-Host "`n[3/3] Generating APK (Debug Mode)..." -ForegroundColor Green
& flutter build apk --debug
if ($LASTEXITCODE -ne 0) {
    Write-Error "Flutter build failed."
    exit 1
}

$apkPath = "build/app/outputs/flutter-apk/app-debug.apk"
if (Test-Path $apkPath) {
    Write-Host "`n✨ SUCCESS! APK generated at: $apkPath" -ForegroundColor Magenta
    Write-Host "Copy this file to your Android device to test." -ForegroundColor White
} else {
    Write-Error "APK not found at expected location."
}

Set-Location -Path "../../.."
