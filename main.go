package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type TorRelay struct {
	Address string
}


type TorBrowserSimulator struct {
	Relays []*TorRelay
}

func NewTorBrowserSimulator(relays []*TorRelay) *TorBrowserSimulator {
	return &TorBrowserSimulator{
		Relays: relays,
	}
}


func (tbs *TorBrowserSimulator) sendRequestThroughRelay(relay *TorRelay, targetURL string) ([]byte, error) {
	
	// proxyURL, err := url.Parse("socks5://localhost:9050")
	proxyURL, err := url.Parse(relay.Address)
	

	if err != nil {
		return nil, err
	}

	
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: time.Second * 10,
	}

	fmt.Println("Sending request through Tor relay...")
	resp, err := client.Get(targetURL)
	if err != nil {
		fmt.Println("Error connecting to Tor relay:", err)
		return nil, err
	}
	fmt.Println("Request successful.")

	defer resp.Body.Close()

	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	// Set up Tor relays.
	relay1 := &TorRelay{Address: "socks5://localhost:9050"}
	
	relays := []*TorRelay{relay1}

	
	tbs := NewTorBrowserSimulator(relays)

	targetURL := "https://www.nasa.gov/"

	
	selectedRelay := relays[0] 
	response, err := tbs.sendRequestThroughRelay(selectedRelay, targetURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the response.
	fmt.Println("Response from Tor relay:", string(response))
}
