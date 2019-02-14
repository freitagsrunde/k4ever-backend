package cmd

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/freitagsrunde/k4ever-backend/internal/context"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Handles user stuff",
	Long:  `This can create delete and update users`,
}

var userCreateCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Creates a user",
	Long: `This will create a new user and will ask
interactively for a password if not given otherwise.
If no DisplayName is given, the username is used instead`,
	Args: cobra.ExactArgs(1),
	Run:  createUser,
}

var password string
var displayName string

func init() {
	userCreateCmd.Flags().StringVarP(&password, "password", "p", "", "Password flag for automated usage")
	userCreateCmd.Flags().StringVarP(&displayName, "displayName", "d", "", "The name that others will see")
	userCmd.AddCommand(userCreateCmd)
	rootCmd.AddCommand(userCmd)
}

func createUser(cmd *cobra.Command, args []string) {
	username := args[0]
	config := context.NewConfig()
	if _, err := k4ever.GetUser(username, config); err != nil {
		if strings.HasPrefix(err.Error(), "record not found") == false {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("User already exists")
		return
	}
	if displayName == "" {
		displayName = args[0]
	}
	if password == "" {
		fmt.Println("Please enter a password")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Could not read password")
			return
		}
		password = string(bytePassword)
		fmt.Println("Please reenter your password")
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Could not read password")
			return
		}
		if password != string(bytePassword) {
			fmt.Println("Passwords don't match")
			return
		}
	}
	fmt.Printf("Create user: %s...\n", username)

	user := models.User{UserName: username, DisplayName: displayName, Password: password}
	if err := k4ever.CreateUser(&user, config); err != nil {
		fmt.Printf("Error while creating user: %s", err.Error())
		return
	}

	fmt.Printf("User %s created\n", username)
}
