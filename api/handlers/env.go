// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package handlers

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/swinslow/obsidian-api/internal/datastore"
)

// Env is the environment for the web handlers.
type Env struct {
	db           datastore.Datastore
	jwtSecretKey string
	oauthConf    *oauth2.Config
	oauthState   string
}

// SetupEnv sets up systems (such as the data store) and variables
// (such as the JWT signing key) that are used across web requests.
func SetupEnv() (*Env, error) {
	// set up datastore
	db, err := datastore.NewDB("host=db sslmode=disable dbname=dev user=postgres-dev")
	if err != nil {
		return nil, err
	}

	err = datastore.InitNewDB(db)
	if err != nil {
		return nil, err
	}

	// set up JWT secret key (from environment)
	JWTSECRETKEY := os.Getenv("JWTSECRETKEY")
	if JWTSECRETKEY == "" {
		return nil, fmt.Errorf("No JWT secret key found; set environment variable JWTSECRETKEY before starting")
	}

	// set up client ID and client secret (from environment)
	GITHUBCLIENTID := os.Getenv("GITHUBCLIENTID")
	if GITHUBCLIENTID == "" {
		return nil, fmt.Errorf("No GitHub client ID found; set environment variable GITHUBCLIENTID before starting")
	}
	GITHUBCLIENTSECRET := os.Getenv("GITHUBCLIENTSECRET")
	if GITHUBCLIENTSECRET == "" {
		return nil, fmt.Errorf("No GitHub client secret found; set environment variable GITHUBCLIENTSECRET before starting")
	}
	OAUTHSTATE := os.Getenv("OAUTHSTATE")
	if OAUTHSTATE == "" {
		return nil, fmt.Errorf("No OAuth state string found; set environment variable OAUTHSTATE before starting")
	}

	oauthConf := &oauth2.Config{
		ClientID:     GITHUBCLIENTID,
		ClientSecret: GITHUBCLIENTSECRET,
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}

	env := &Env{
		db:           db,
		jwtSecretKey: JWTSECRETKEY,
		oauthConf: oauthConf,
		oauthState: OAUTHSTATE,
	}
	return env, nil
}

// getTestEnv creates the Env object used for the handlers
// unit test suite. It is not exported and should NEVER be
// called by production code.
func getTestEnv() *Env {
	db := &mockDB{}

	oauthConf := &oauth2.Config{
		ClientID:     "abcdef0123abcdef4567",
		ClientSecret: "abcdef0123abcdef4567abcdef8901abcdef2345",
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}

	env := &Env{
		db:           db,
		jwtSecretKey: "keyForTesting",
		oauthConf:    oauthConf,
		oauthState:   "nonRandomStateString",
	}
	return env
}