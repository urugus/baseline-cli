package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urugus/baseline-cli/internal/config"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage local configuration",
	}

	cmd.AddCommand(newConfigPathCommand())
	cmd.AddCommand(newConfigGetCommand())
	cmd.AddCommand(newConfigSetCommand())
	cmd.AddCommand(newConfigUnsetCommand())
	return cmd
}

func newConfigPathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print config path",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.Path()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get api-key",
		Short: "Print masked API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "api-key" {
				return fmt.Errorf("unsupported config key: %s", args[0])
			}
			value, source, err := config.APIKey()
			if err != nil {
				return err
			}
			if value == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "api-key is not set")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "api-key=%s source=%s\n", config.Mask(value), source)
			return nil
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set api-key [value]",
		Short: "Set API key",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "api-key" {
				return fmt.Errorf("unsupported config key: %s", args[0])
			}

			value := ""
			if len(args) == 2 {
				value = args[1]
			} else {
				fmt.Fprint(cmd.ErrOrStderr(), "API key: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					value = scanner.Text()
				}
				if err := scanner.Err(); err != nil {
					return err
				}
			}

			if err := config.SetAPIKey(strings.TrimSpace(value)); err != nil {
				return err
			}
			path, err := config.Path()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "api-key saved to %s\n", path)
			return nil
		},
	}
}

func newConfigUnsetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unset api-key",
		Short: "Unset API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "api-key" {
				return fmt.Errorf("unsupported config key: %s", args[0])
			}
			if err := config.UnsetAPIKey(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "api-key unset")
			return nil
		},
	}
}
