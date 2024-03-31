package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	PORT         = "6942"
	wg           sync.WaitGroup
	currentToken *oauth2.Token
)

func StartServer(tokenChan chan *oauth2.Token) {
	server := &http.Server{Addr: ":" + PORT}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Authorization successful. You can close this window.")
		token, err := getConfig().Exchange(context.Background(), code)
		if err != nil {
			log.Printf("Failed to exchange code for token: %v", err)
			return
		}
		tokenChan <- token // Send the token to the channel
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case token := <-tokenChan: // Corrected: Now able to receive from tokenChan
			currentToken = token // Store the token in the global variable
		case <-time.After(5 * time.Minute):
			log.Println("Authorization timed out")
		}
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown server: %v", err)
		}
	}()
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:" + PORT,
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.readonly"},
	}
}

func RefreshToken(config *oauth2.Config) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if currentToken == nil {
			log.Println("No token available for refresh.")
			continue
		}

		tokenSource := config.TokenSource(context.Background(), currentToken)
		newToken, err := tokenSource.Token()
		if err != nil {
			log.Printf("Error refreshing token: %v\n", err)
			continue
		}

		if newToken.AccessToken != currentToken.AccessToken {
			log.Println("Token refreshed")
			currentToken = newToken // Update the global token
		}
	}
}

func Authenticate() error {
	tokenChan := make(chan *oauth2.Token)
	StartServer(tokenChan)

	url := getConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Println("Open the following URL in your browser to authenticate:", url)

	err := exec.Command("xdg-open", url).Run()
	if err != nil {
		return err
	}

	wg.Wait() // Wait for the authentication process to complete

	go RefreshToken(getConfig())

	return nil
}
