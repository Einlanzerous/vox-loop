package main

import (
	"fmt"
	"os"

	"github.com/Einlanzerous/vox-loop/internal/admin"
	"github.com/Einlanzerous/vox-loop/internal/preflight"
	"github.com/Einlanzerous/vox-loop/internal/setup"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "vox-loop",
		Short:   "Vox Loop — Matrix/Dendrite homeserver manager",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return preflight.RunAndExec(setup.GenerateContainerConfig)
		},
		SilenceUsage: true,
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Generate Dendrite config, keys, and well-known files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setup.Run()
		},
	}

	createAccountCmd := &cobra.Command{
		Use:   "create-account",
		Short: "Create a new Matrix user account",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			isAdmin, _ := cmd.Flags().GetBool("admin")
			if username == "" {
				return fmt.Errorf("--username is required")
			}
			return admin.CreateAccount(username, isAdmin)
		},
	}
	createAccountCmd.Flags().String("username", "", "Username for the new account")
	createAccountCmd.Flags().Bool("admin", false, "Grant admin privileges")

	adminCmd := &cobra.Command{
		Use:   "admin",
		Short: "Administrative commands",
	}
	adminCmd.AddCommand(createAccountCmd)

	rootCmd.AddCommand(initCmd, adminCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
