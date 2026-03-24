//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func checkEndpoint(name, url string) {
	client := &http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	resp, err := client.Get(url)
	latency := time.Since(start)

	if err != nil {
		fmt.Printf("[❌] %s (%s): FALHOU (%v) - Erro: %v\n", name, url, latency, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[✅] %s (%s): ONLINE (%v) - Status: %s\n", name, url, latency, resp.Status)
}

func checkPort(name, host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		fmt.Printf("[❌] %s (%s): FECHADA/BLOQUEADA - Erro: %v\n", name, addr, err)
		return
	}
	conn.Close()
	fmt.Printf("[✅] %s (%s): ABERTA\n", name, addr)
}

func checkExec(name, cmdName string) {
	path, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Printf("[❌] %s (%s): NÃO ENCONTRADO NO PATH\n", name, cmdName)
		return
	}
	fmt.Printf("[✅] %s (%s): LOCALIZADO EM %s\n", name, cmdName, path)
}

func main() {
	fmt.Println("--- AUDITORIA DE SISTEMA STARDF-ANIME ---")

	// 1. APIs Externas
	checkEndpoint("AniList API", "https://graphql.anilist.co")

	// 2. Serviços Locais
	checkPort("Spider Node.js Bridge", "localhost", 3000)
	checkPort("StarDF-Anime WebUI", "localhost", 8080)

	// 3. Executáveis Críticos
	checkExec("Video Player", "mpv")
	checkExec("Streamlink Service", "streamlink")

	// 4. Acesso a Pastas Locais
	dir, err := os.Stat("internal/watchlist")
	if err != nil || !dir.IsDir() {
		fmt.Println("[❌] Persistência (internal/watchlist): DIRETÓRIO AUSENTE OU INACESSÍVEL")
	} else {
		fmt.Println("[✅] Persistência (internal/watchlist): OK")
	}
}
