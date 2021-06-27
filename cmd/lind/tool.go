// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lind

import (
	"fmt"
	"os"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

func printLogoWhenIsTty() {
	if logger.IsTerminal(os.Stdout) {
		_, _ = fmt.Fprintf(os.Stdout, logger.Cyan.Add(linDBLogo))
		_, _ = fmt.Fprintf(os.Stdout, logger.Green.Add(" ::  LinDB  :: ")+
			fmt.Sprintf("%22s", fmt.Sprintf("(v%s Release)", getVersion())))
		_, _ = fmt.Fprintf(os.Stdout, "\n\n")
	}
}

// askForConfirmation use Scanln to parse user input.
// A user must type in "yes" or "no" and then press enter.
// It has fuzzy matching, so "y", "Y", "yes", "YES"
func askForConfirmation() (confirmed bool, err error) {
	_, _ = fmt.Fprintf(os.Stdout, "[y/n]:")

	var answer string
	if _, err = fmt.Scanln(&answer); err != nil {
		return false, err
	}
	okayAnswers := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayAnswers := []string{"n", "N", "no", "No", "NO"}

	containsString := func(all []string, element string) bool {
		for _, item := range all {
			if item == element {
				return true
			}
		}
		return false
	}
	if containsString(okayAnswers, answer) {
		return true, nil
	} else if containsString(nokayAnswers, answer) {
		return false, nil
	} else {
		return askForConfirmation()
	}
}

func checkExistenceOf(path string) error {
	if fileutil.Exist(path) {
		fmt.Println("An old config already exist, do you want to overwrite it?")
		confirmed, err := askForConfirmation()
		if err != nil {
			return err
		}
		if confirmed {
			fmt.Println("Overwriting...")
		} else {
			fmt.Println("Skipping...")
			os.Exit(0)
		}
	}
	return nil
}
