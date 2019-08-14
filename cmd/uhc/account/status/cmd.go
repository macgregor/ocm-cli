/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package status

import (
	"fmt"

	"github.com/spf13/cobra"

	acc_util "github.com/openshift-online/uhc-cli/pkg/account"
	"github.com/openshift-online/uhc-cli/pkg/config"
)

var args struct {
	debug bool
}

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Status of current user.",
	Long:  "Display status of current user.",
	RunE:  run,
}

func init() {
	// Add flags to rootCmd:
	flags := Cmd.Flags()
	flags.BoolVar(
		&args.debug,
		"debug",
		false,
		"Enable debug mode.",
	)
}

func run(cmd *cobra.Command, argv []string) error {

	// Load the configuration file:
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Can't load config file: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("Not logged in, run the 'login' command")
	}

	// Check that the configuration has credentials or tokens that haven't have expired:
	armed, err := cfg.Armed()
	if err != nil {
		return fmt.Errorf("Can't check if tokens have expired: %v", err)
	}
	if !armed {
		return fmt.Errorf("Tokens have expired, run the 'login' command")
	}

	// Create the connection, and remember to close it:
	connection, err := cfg.Connection()
	if err != nil {
		return fmt.Errorf("Can't create connection: %v", err)
	}
	defer connection.Close()

	// Send the request:
	response, err := connection.AccountsMgmt().V1().CurrentAccount().Get().
		Send()
	if err != nil {
		return fmt.Errorf("Can't get current account: %v", err)
	}

	// Display user and which server they are logged into
	currAccount := response.Body()
	fmt.Println(fmt.Sprintf("%s on %s", currAccount.Username(), cfg.URL))

	// Display roles currently assigned to the user
	roleSlice, err := acc_util.GetRolesFromUser(currAccount, connection)
	if err != nil {
		return err
	}
	fmt.Printf("Roles: %v\n", roleSlice)

	return nil
}
