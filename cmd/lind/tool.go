package lind

import (
	"fmt"
	"os"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

func printLogoWhenIsTty() {
	if logger.IsTerminal(os.Stdout) {
		fmt.Fprintf(os.Stdout, logger.Cyan.Add(linDBLogo))
		fmt.Fprintf(os.Stdout, logger.Green.Add(" ::  LinDB  :: ")+
			fmt.Sprintf("%22s", fmt.Sprintf("(v%s Release)", getVersion())))
		fmt.Fprintf(os.Stdout, "\n\n")
	}
}

// askForConfirmation use Scanln to parse user input.
// A user must type in "yes" or "no" and then press enter.
// It has fuzzy matching, so "y", "Y", "yes", "YES"
func askForConfirmation() (confirmed bool, err error) {
	fmt.Fprintf(os.Stdout, "[y/n]:")

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
