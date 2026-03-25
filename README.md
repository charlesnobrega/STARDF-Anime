<h4 align="center">
    <p>
        <b>English</b> |
        <a href="https://github.com/charlesnobrega/STARDF-Anime/blob/main/README_pt-br.md">Рortuguês</a>
    </p>
</h4>

<p align="center">
  <img src="docs/banner.png" alt="StarDF-Anime Banner" width="800"/>
</p>

[![GitHub license](https://img.shields.io/github/license/charlesnobrega/STARDF-Anime)](https://github.com/charlesnobrega/STARDF-Anime/blob/main/LICENSE)
![GitHub stars](https://img.shields.io/github/stars/charlesnobrega/STARDF-Anime)
![GitHub last commit](https://img.shields.io/github/last-commit/charlesnobrega/STARDF-Anime)
![GitHub forks](https://img.shields.io/github/forks/charlesnobrega/STARDF-Anime?style=social)
[![Build Status](https://github.com/charlesnobrega/STARDF-Anime/actions/workflows/ci.yml/badge.svg)](https://github.com/charlesnobrega/STARDF-Anime/actions)
![GitHub contributors](https://img.shields.io/github/contributors/charlesnobrega/STARDF-Anime)

# StarDF-Anime

StarDF-Anime is a powerful terminal user interface (TUI) for browsing, streaming, and tracking anime & movies. It features real-time synchronization with AniList, high-performance scraping, and metadata enrichment for Portuguese and English content.

### Mobile Version (Em breve)

A mobile version of StarDF-Anime is planned for Android devices.

### Community

Join our official community for support, updates, and feedback:
[![Discord](https://img.shields.io/discord/1234567890?color=5865F2&label=Discord&logo=discord&logoColor=white)](https://discord.gg/stardf-anime)

## Features

- Search for anime by name
- Browse episodes
- Support subbed and dubbed content in English and Portuguese
- Skip anime Intro
- Play online with quality selection
- Download single episodes
- Discord RPC about the anime
- Batch download multiple episodes
- **NEW:** Premium Web User Interface (same binary as CLI, via `-web`)
- **NEW:** Movies and TV Shows support via FlixHQ source
- **NEW:** OMDb integration for movie/TV metadata (ratings, genres, runtime)
- **NEW:** Universal SQLite tracking (100% Go, works on all platforms without CGO)

> **Note:** StarDF-Anime now uses a pure Go SQLite implementation. All official release binaries include full tracking and watch history support by default, without requiring external CGO dependencies.

# Web UI Demo

![StarDF-Anime Web UI Demo](docs/web_demo.webp)

## Prerequisites

- Go (at latest version)

- Mpv(at latest version)

## how to install and run

### Universal install (Only needs Go installed)  

```shell
go install github.com/charlesnobrega/STARDF-Anime/cmd/stardf-anime@latest
```

...existing code...
### Manual install methods

```shell
git clone https://github.com/charlesnobrega/STARDF-Anime.git
```

```shell
cd STARDF-Anime
```

```shell
go run cmd/stardf-anime/main.go
```

## Movies and TV Shows

StarDF-Anime supports movies and TV shows through the FlixHQ source. Use the `--source flixhq` flag to search for movies and TV shows. You can also restrict results by type using the `--type` parameter (for example `--type movie` to search only movies).

```bash
# Search for movies/TV shows
stardf-anime --source flixhq "The Matrix"

# Search specifically for movies
stardf-anime --source flixhq --type movie "Inception"

# Search specifically for TV shows
stardf-anime --source flixhq --type tv "Breaking Bad"

# Enable subtitles (English by default)
stardf-anime --source flixhq --subs "Avatar"
```


## Linux

<details>
<summary>Arch Linux / Manjaro (AUR-based systems)</summary>

Using Yay:

```bash
yay -S stardf-anime
```

or using Paru:

```bash
paru -S stardf-anime
```

Or, to manually clone and install:

```bash
git clone https://aur.archlinux.org/stardf-anime.git
cd stardf-anime
makepkg -si
sudo pacman -S mpv
```

</details>

<details>
<summary>Debian / Ubuntu (and derivatives)</summary>

```bash
sudo apt update
sudo apt install mpv

# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux-amd64
```

</details>

<details>
<summary>Fedora Installation</summary>

```bash
sudo dnf update
sudo dnf install mpv

# For x86_64 systems (Em breve):
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux
```

</details>

<details>
<summary>openSUSE Installation</summary>

```bash
sudo zypper refresh
sudo zypper install mpv

# For x86_64 systems (Em breve):
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux
```

</details>

## Windows

Option 1: Windows Executable (single binary: CLI + Web UI)

- Download the latest `stardf-anime-windows.zip` from the [releases](https://github.com/charlesnobrega/STARDF-Anime/releases) section.
- Extract and run `stardf-anime.exe`.
- Use `stardf-anime.exe -web` to launch the Premium Web UI.

Option 2: Using the installer

- An Inno Setup based installer is available for easier integration (PATH + icon + shortcuts).
- It creates Start Menu shortcuts for `StarDF-Anime (Console)` and `StarDF-Anime (Web UI)`.

## macOS

<details>
<summary>macOS Installation</summary>

First, install mpv using Homebrew:

```bash
# Install Homebrew if you haven't already
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install mpv
brew install mpv

# Download and install stardf-anime
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-darwin-arm64
```

Alternative installation using MacPorts:

```bash
# Install mpv using MacPorts
sudo port install mpv

# Download and install StarDF-Anime (Em breve)
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-apple-darwin
```

</details>

### Usage

To start the application, simply run:

```bash
stardf-anime
```

### Advanced Usage

You can search for and play content directly from the command line:

- Search and play:
```bash
stardf-anime "One Piece"
```

- Update to latest version:
```bash
stardf-anime --update
```

- Help and options:
```bash
stardf-anime --help
```

The program provides a fully interactive TUI. You can navigate through search results, selection screens, and playback controls using your keyboard.

