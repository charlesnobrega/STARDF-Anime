package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/alvarorichard/Goanime/internal/util"
)

func main() {
	url := "https://animesonlinecc.to/?s=One+Piece"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", util.UserAgentList())
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
