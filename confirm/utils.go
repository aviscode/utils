package confirm

import (
	"fmt"
	"log"
	"strings"
)

// AskForConfirmation  uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// the func gets the question to ask if the value in empty the default question is: WARNING: Are you sure? (yes/no)
func AskForConfirmation(question string) bool {
	if question == "" {
		fmt.Print("WARNING: Are you sure? (yes/no): ")
	} else {
		fmt.Println(question)
	}
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return AskForConfirmation(question)
	}
}
