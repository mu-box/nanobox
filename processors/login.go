package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/odin"

	printutil "github.com/sdomino/go-util/print"
)

// Process ...
func Login(username, password string) error {

	// request Username/Password if missing
	if username == "" {
		// add in tylers display system for prompting
		username = printutil.Prompt("Username: ")
	}

	if password == "" {
		// ReadPassword prints Password: already
		pass, err := util.ReadPassword()
		if err != nil {
			// TODO: print out the error to the log
			return fmt.Errorf("failed to read password: %s", err.Error())
		}
		password = pass
	}

	// verify that the user exists
	token, err := odin.Auth(username, password)
	if err != nil {
		return fmt.Errorf("unable to authenticate with nanobox: %s", err.Error())
	}

	// store the user token
	auth := models.Auth{Key: token}
	if auth.Save() != nil {
		return fmt.Errorf("unable to save user")
	}

	fmt.Println("TODO: Message: user has been verified and granted access")

	return nil
}
