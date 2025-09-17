package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/lfr-tools/internal/keychain"
)

var keychainCmd = &cobra.Command{
	Use:   "keychain",
	Short: "Manage secure token storage in system keychain",
	Long:  `Store and manage access tokens securely using platform-native credential storage.`,
}

var keychainStoreCmd = &cobra.Command{
	Use:   "store [project] [username] [token]",
	Short: "Store token in system keychain",
	Long: `Store an access token securely in the system keychain (macOS Keychain,
Windows Credential Manager, or Linux secret service).`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		username := args[1]
		token := args[2]

		return storeTokenInKeychain(project, username, token)
	},
}

var keychainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored tokens",
	Long:  `List all access tokens stored in the system keychain.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listKeychainTokens()
	},
}

var keychainDeleteCmd = &cobra.Command{
	Use:   "delete [project] [username]",
	Short: "Delete token from keychain",
	Long:  `Remove an access token from the system keychain.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := args[0]
		username := args[1]

		return deleteTokenFromKeychain(project, username)
	},
}

func init() {
	rootCmd.AddCommand(keychainCmd)

	keychainCmd.AddCommand(keychainStoreCmd)
	keychainCmd.AddCommand(keychainListCmd)
	keychainCmd.AddCommand(keychainDeleteCmd)
}

// storeTokenInKeychain stores a token securely.
func storeTokenInKeychain(project, username, token string) error {
	tokenStore, err := keychain.NewTokenStore()
	if err != nil {
		return fmt.Errorf("failed to initialize token store: %w", err)
	}

	err = tokenStore.StoreToken(project, username, token)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	fmt.Printf("✅ Token stored securely for %s in project %s\n", username, project)
	fmt.Printf("Platform: %s keychain integration\n", getPlatformName())

	return nil
}

// listKeychainTokens lists all stored tokens.
func listKeychainTokens() error {
	tokenStore, err := keychain.NewTokenStore()
	if err != nil {
		return fmt.Errorf("failed to initialize token store: %w", err)
	}

	accounts, err := tokenStore.ListTokens()
	if err != nil {
		return fmt.Errorf("failed to list tokens: %w", err)
	}

	if len(accounts) == 0 {
		fmt.Printf("No tokens stored in keychain.\n")
		return nil
	}

	fmt.Printf("Stored tokens in %s keychain:\n\n", getPlatformName())
	fmt.Printf("%-20s\n", "ACCOUNT")
	fmt.Println(strings.Repeat("-", 25))

	for _, account := range accounts {
		fmt.Printf("%-20s\n", account)
	}

	fmt.Printf("\nTotal: %d tokens\n", len(accounts))
	return nil
}

// deleteTokenFromKeychain removes a token.
func deleteTokenFromKeychain(project, username string) error {
	tokenStore, err := keychain.NewTokenStore()
	if err != nil {
		return fmt.Errorf("failed to initialize token store: %w", err)
	}

	err = tokenStore.DeleteToken(project, username)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	fmt.Printf("✅ Token deleted for %s in project %s\n", username, project)

	return nil
}

// getPlatformName returns a user-friendly platform name.
func getPlatformName() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	default:
		return "File-based"
	}
}