/*
 * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package googledrive holds the implementation related to Drive interactions.
package googledrive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	// FileMimeType represents the GDrive file mime type.
	FileMimeType = "mimeType = 'application/vnd.googledrive-apps.file'"
	// DirMimeType represents the GDrive directory mime type.
	DirMimeType = "mimeType = 'application/vnd.googledrive-apps.folder'"
)

var (
	// ClientSecret is the client secret for the GDrive client and must be set from the LD flags.
	ClientSecret string
	// ClientID is the client ID for the GDrive client and must be set from the LD flags.
	ClientID string
	creds    = &credentials{
		Installed: installedCreds{
			RedirectURIs: []string{"urn:ietf:wg:oauth:2.0:oob", "http://localhost"},
			AuthURI:      "https://accounts.google.com/o/oauth2/auth",
			TokenURI:     "https://oauth2.googleapis.com/token",
			ClientSecret: ClientSecret,
			ClientID:     ClientID,
		},
	}
)

// credentials struct represents the credentials required for the Google Drive client.
type credentials struct {
	Installed installedCreds `json:"installed"`
}

type installedCreds struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}

// Client represents the GDrive client.
type Client struct {
	srv       *drive.Service
	TokenFile string
}

// Retrieve a token, saves the token, then returns the generated client.
func (d *Client) getClient(config *oauth2.Config) (*http.Client, error) {
	// The file TokenFile stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(d.TokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		err := saveToken(d.TokenFile, tok)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "Unable to cache oauth token")
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	return json.NewEncoder(f).Encode(token)
}

// Init initialises the GDrive client.
func (d *Client) Init() error {
	b, err := utils.MarshalData(creds)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal creds")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return errors.Wrap(err, "unable to get configurations")
	}

	client, err := d.getClient(config)
	if err != nil {
		return errors.Wrap(err, "unable to parse client secret file to config")
	}
	d.srv, err = drive.New(client)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve Drive client")
	}
	return nil
}

// IsDirExists method checks whether the given directory exists or not.
func (d *Client) IsDirExists(dirName string) (bool, string, error) {
	dirList, err := d.srv.Files.List().Q(DirMimeType).Q("name = '" + dirName + "'").Do()
	if err != nil {
		return false, "", err
	}
	c := len(dirList.Files)
	if c == 0 {
		return false, "", nil
	} else if c != 1 {
		return false, "", errors.New(fmt.Sprintf("More than one directory exists under the given name: %s", dirName))
	}
	return true, dirList.Files[0].Id, nil
}

// IsFileExists method checks whether the given file exists under the given directory or not.
func (d *Client) IsFileExists(dirID string, fileName string) (bool, string, error) {
	fileSearchQuery := "name = '" + fileName + "' and '" + dirID + "' in parents"
	fileList, err := d.srv.Files.List().Q(FileMimeType).Q(fileSearchQuery).Do()
	if err != nil {
		return false, "", err
	}
	c := len(fileList.Files)
	if c == 0 {
		return false, "", nil
	} else if c != 1 {
		return false, "", errors.New(fmt.Sprintf("More than one file exists under the given name: %s", fileName))
	}
	return true, fileList.Files[0].Id, nil
}

// FileContent method downloads the given file content.
func (d *Client) FileContent(fileID string) ([]byte, error) {
	resp, err := d.srv.Files.Get(fileID).Download()
	if err != nil {
		return nil, errors.Wrap(err, "unable to download password file")
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			//todo handle
			panic(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read password DB file from of the response")
	}
	return body, nil
}

// CreateFile method creates the given file.
func (d *Client) CreateFile(name string, mimeType string, content io.Reader, parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}
	file, err := d.srv.Files.Create(f).Media(content).Do()

	if err != nil {
		return nil, err
	}
	return file, nil
}

// CopyFile method copies the given file.
func (d *Client) CopyFile(fileID, name, mimeType string, parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}
	file, err := d.srv.Files.Copy(fileID, f).Do()
	if err != nil {
		return nil, err
	}
	return file, nil
}

// CreateDir method creates the given directory under the given parent ID.
func (d *Client) CreateDir(name string, parentID string) (*drive.File, error) {
	dir := &drive.File{
		Name:     name,
		MimeType: DirMimeType,
		Parents:  []string{parentID},
	}
	file, err := d.srv.Files.Create(dir).Do()
	if err != nil {
		return nil, err
	}
	return file, nil
}

// UpdateFileContent method updates the given file with given content.
func (d *Client) UpdateFileContent(fileID string, content []byte) error {
	_, err := d.srv.Files.Update(fileID, &drive.File{}).Media(bytes.NewBuffer(content)).Do()
	if err != nil {
		return errors.Wrap(err, "unable to update file")
	}
	return nil
}
