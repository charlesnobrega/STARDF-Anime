package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "https://superanimes.in/busca/?search_query=naruto"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	
	f, _ := os.Create("superanimes_naruto.html")
	defer f.Close()
	io.Copy(f, resp.Body)
	fmt.Println("Saved to superanimes_naruto.html")
}
