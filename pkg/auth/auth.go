package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	gAuthFilePath = os.Getenv("GOOGLE_CREDS")
	tokenFilePath = os.Getenv("GOOGLE_TOKEN")
)

func NewGmailClient() (*http.Client, error) {
	if gAuthFilePath == "" {
		return nil, fmt.Errorf("no path to credentail.json provided")
	}

	if tokenFilePath == "" {
		return nil, fmt.Errorf("no path to token.json provided")
	}

	b, err := ioutil.ReadFile(gAuthFilePath)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		return nil, err
	}

	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {

		if tok, err = getTokenFromWeb(config); err != nil {
			return nil, err
		}

		if err = saveToken(gAuthFilePath, tok); err != nil {
			return nil, err
		}
	}

	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	log.Info("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
