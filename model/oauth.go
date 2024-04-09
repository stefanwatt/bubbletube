package model

import (
	"bubbletube/config"
	"context"
	"encoding/json"
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
		token, err := getConfig().Exchange(context.Background(), code)
		if err != nil {
			log.Printf("Failed to exchange code for token: %v", err)
			return
		}
		tokenChan <- token
		SaveToken(token)
		w.Write([]byte("Authorization successful. You can now close this tab and return to the terminal where bubbletube is running."))
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
			currentToken = newToken
			SaveToken(newToken)
		}
	}
}

func Authenticate() error {
	config := getConfig()
	loadedToken, err := LoadToken()
	if err == nil {
		currentToken = loadedToken
	} else {
		fmt.Println("failed to load stored token ", err)
	}
	if currentToken != nil && currentToken.RefreshToken != "" {
		tokenSource := config.TokenSource(context.Background(), currentToken)
		newToken, err := tokenSource.Token()

		if err != nil {
			log.Printf("Failed to refresh token: %v", err)
		} else {
			currentToken = newToken
			return nil
		}
	} else {
		if currentToken == nil {
			fmt.Println("NOT Refreshing token because currentToken is nil")
		} else if currentToken.RefreshToken == "" {
			fmt.Println("NOT Refreshing token because currentToken.RefreshToken is empty")
		}
	}
	tokenChan := make(chan *oauth2.Token)
	StartServer(tokenChan)
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	err = exec.Command("xdg-open", url).Run()
	if err != nil {
		return err
	}
	wg.Wait()
	go RefreshToken(config)
	return nil
}

func SaveToken(token *oauth2.Token) error {
	file, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		fmt.Println("Error marshalling token:", err)
		return err
	}
	return os.WriteFile(config.TOKEN_PATH, file, 0600)
}

func LoadToken() (*oauth2.Token, error) {
	file, err := os.ReadFile(config.TOKEN_PATH)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(file, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
