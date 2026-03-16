<h4 align="center">
    <p>
        <b>Рortuguês</b> |
        <a href="https://github.com/charlesnobrega/STARDF-Anime/blob/main/README.md">English</a>
    </p>
</h4>

<p align="center">
  <img src="docs/logo_stardf.png" alt="StarDF-Anime Logo" width="400"/>
</p>

[![GitHub license](https://img.shields.io/github/license/charlesnobrega/STARDF-Anime)](https://github.com/charlesnobrega/STARDF-Anime/blob/main/LICENSE)
![GitHub stars](https://img.shields.io/github/stars/charlesnobrega/STARDF-Anime)
![GitHub last commit](https://img.shields.io/github/last-commit/charlesnobrega/STARDF-Anime)
![GitHub forks](https://img.shields.io/github/forks/charlesnobrega/STARDF-Anime?style=social)
[![Build Status](https://github.com/charlesnobrega/STARDF-Anime/actions/workflows/ci.yml/badge.svg)](https://github.com/charlesnobrega/STARDF-Anime/actions)
![GitHub contributors](https://img.shields.io/github/contributors/charlesnobrega/STARDF-Anime)

# StarDF-Anime

StarDF-Anime é uma interface de usuário baseada em texto (TUI) poderosa, desenvolvida em Go, evoluindo do StarDF-Anime original. Ele permite aos usuários procurar animes, filmes e séries, e reproduzir ou baixar conteúdos diretamente no mpv. É especificamente otimizado para scraping de alta performance e enriquecimento de metadados para conteúdos em português e inglês.

### Versão Mobile (Em breve)

Uma versão mobile do StarDF-Anime está planejada para dispositivos Android.

### Comunidade (Em breve)

Um servidor oficial no Discord está nos planos para suporte e novidades.

## Recursos

- Buscar anime por nome
- Navegar pelos episódios
- Suporte a conteúdo legendado e dublado em inglês e português
- Pular introdução do anime
- Reproduzir online com seleção de qualidade
- Baixar episódios únicos
- Discord RPC sobre o anime
- Download em lote de múltiplos episódios
- Retomar reprodução de onde parou (em builds com suporte SQLite)
- Rastrear episódios assistidos (em builds com suporte SQLite)

> **Nota:** StarDF-Anime pode ser compilado com ou sem suporte SQLite para rastreamento do progresso do anime.  
> [Veja a documentação das opções de build](docs/BUILD_OPTIONS.md) para mais detalhes.

> ⚠️ Aviso: disponibilidade da fonte em Português (PT-BR)
>
> Se você deseja assistir animes em português (PT-BR) e está fora do Brasil, será necessário usar uma VPN, proxy ou qualquer método para obter um endereço de IP brasileiro. A fonte de animes em PT-BR bloqueia o acesso de IPs fora do Brasil.

# Demo

<https://github.com/charlesnobrega/STARDF-Anime/assets/88117897/ffec6ad7-6ac1-464d-b048-c80082119836>

## Pré-requisitos

- Go (na versão mais recente)
- Mpv (na versão mais recente)

## Como instalar e executar

### Instalação Universal (Só precisa do Go instalado)

```shell
go install github.com/charlesnobrega/STARDF-Anime/cmd/stardf-anime@latest
```

### Métodos de instalação manual

```shell
git clone https://github.com/charlesnobrega/STARDF-Anime.git
```

```shell
cd STARDF-Anime
```

```shell
go run cmd/stardf-anime/main.go
```

## Filmes e Séries

StarDF-Anime agora suporta filmes e séries através da fonte FlixHQ. Use a flag `--source flixhq` para buscar filmes e séries. Você também pode filtrar por tipo usando o parâmetro `--type` (por exemplo `--type movie` para buscar somente filmes).

```bash
# Buscar filmes/séries
stardf-anime --source flixhq "The Matrix"

# Buscar somente filmes
stardf-anime --source flixhq --type movie "Inception"

# Buscar somente séries
stardf-anime --source flixhq --type tv "Breaking Bad"

# Habilitar legendas (inglês por padrão)
stardf-anime --source flixhq --subs "Avatar"
```



## Linux

<details>
<summary>Arch Linux / Manjaro (sistemas baseados em AUR)</summary>

Usando Yay:

```bash
yay -S StarDF-Anime
```

ou usando Paru:

```bash
paru -S StarDF-Anime
```

Ou, para clonar e instalar manualmente:

```bash
git clone https://aur.archlinux.org/StarDF-Anime.git
cd StarDF-Anime
makepkg -si
sudo pacman -S mpv
```

</details>

<details>
<summary>Debian / Ubuntu (e derivados)</summary>

```bash
sudo apt update
sudo apt install mpv

# Para sistemas x86_64 (Em breve):
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux
```

</details>

<details>
<summary>Instalação no Fedora</summary>

```bash
sudo dnf update
sudo dnf install mpv

# Para sistemas x86_64 (Em breve):
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux
```

</details>

<details>
<summary>Instalação no openSUSE</summary>

```bash
sudo zypper refresh
sudo zypper install mpv

# Para sistemas x86_64 (Em breve):
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-linux
```

</details>

## Windows

<details>
<summary>Instalação no Windows</summary>

> **Altamente Recomendado:** Use o instalador para a melhor experiência no Windows.

Opção 1: Usando o instalador (Em breve)

- Um instalador Windows estará disponível na próxima versão.

Opção 2: Executável independente (Em breve)

- Executáveis estarão disponíveis na seção de [releases](https://github.com/charlesnobrega/STARDF-Anime/releases) em breve.

</details>

## macOS

<details>
<summary>Instalação no macOS</summary>

Primeiro, instale o mpv usando o Homebrew:

```bash
# Instale o Homebrew se você ainda não tiver
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Instale o mpv
brew install mpv

# Baixe e instale o StarDF-Anime (Em breve)
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-apple-darwin
```

Instalação alternativa usando MacPorts:

```bash
# Instale o mpv usando MacPorts
sudo port install mpv

# Baixe e instale o StarDF-Anime (Em breve)
# curl -Lo stardf-anime https://github.com/charlesnobrega/STARDF-Anime/releases/latest/download/stardf-anime-apple-darwin
```

</details>

### Passos de Configuração Adicionais

# Instalação no NixOS (Flakes)

## Execução Temporária

```shell
nix github:charlesnobrega/STARDF-Anime
```

## Instalação

Adicione no seu `flake.nix`:

```nix
 inputs.StarDF-Anime.url = "github:charlesnobrega/STARDF-Anime";
```

Passe as entradas para seus módulos usando `specialArgs` e então no `configuration.nix`:

```nix
environment.systemPackages = [
  inputs.StarDF-Anime.packages.${pkgs.system}.StarDF-Anime
];
```

### Uso no Linux e macOS

```shell
go-anime
```

### Uso no Windows

```shell
StarDF-Anime
```

### Uso Avançado

Você também pode usar parâmetros para procurar e reproduzir anime diretamente. Aqui estão alguns exemplos:

- Para procurar e reproduzir um anime diretamente, use o seguinte comando:

```shell
StarDF-Anime  "nome do anime"
```

- Para atualizar o StarDF-Anime para a versão mais recente, use a flag de atualização:

```shell
StarDF-Anime --update
```

Este comando irá automaticamente baixar e instalar a versão mais recente do StarDF-Anime usando o mecanismo de atualização integrado do Go.

Você pode usar a opção `-h` ou `--help` para exibir informações de ajuda sobre como usar o comando `StarDF-Anime`.

```shell
StarDF-Anime -h
```

O programa solicitará que você insira o nome de um anime. Digite o nome do anime que deseja assistir.

O programa apresentará uma lista de animes que correspondem à sua entrada. Navegue pela lista usando as setas do teclado e pressione enter para selecionar um anime.

Em seguida, o programa apresentará uma lista de episódios do anime selecionado. Novamente, navegue pela lista usando as setas do teclado e pressione enter para selecionar um episódio.

O episódio selecionado será então reproduzido no MPV.

# Agradecimentos

[@KitsuneSemCalda](https://github.com/KitsuneSemCalda), [@RushikeshGaikwad](https://github.com/Wraient) e [@the-eduardo](https://github.com/the-eduardo) por ajudar e melhorar essa aplicação.

# Alternativas

Se você estiver procurando por mais opções, confira este projeto alternativo do meu amigo [@KitsuneSemCalda](https://github.com/KitsuneSemCalda) chamado [Animatic-v2](https://github.com/KitsuneSemCalda/Animatic-v2), que foi inspirado no StarDF-Anime.

## Contribuindo

Contribuições para melhorar ou aprimorar são sempre bem-vindas. Antes de contribuir, por favor leia nosso guia de desenvolvimento abrangente para informações detalhadas sobre nosso fluxo de trabalho, padrões de código e estrutura do projeto.

📖 **[Guia de Desenvolvimento](docs/Development.md)** - Leitura essencial para contribuidores

**Início Rápido para Contribuidores:**

1. Faça um fork do projeto
2. Leia o [Guia de Desenvolvimento](docs/Development.md) completamente
3. Crie sua branch de funcionalidade a partir de `dev` (nunca de `main`)
4. Siga nossos padrões de código (use `go fmt`, adicione comentários significativos)
5. Certifique-se de que todos os testes passem e adicione testes para novas funcionalidades
6. Faça commit das suas alterações usando formato de commit convencional
7. Faça push para sua branch
8. Abra um Pull Request para a branch `dev`

**Importante:** Nunca faça commit diretamente na branch `main`. Todas as mudanças devem passar pela branch `dev` primeiro.
