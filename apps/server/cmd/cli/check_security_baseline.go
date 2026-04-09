package cli

import (
	"fmt"
	"io"
	"strings"

	appbootstrap "servify/apps/server/internal/app/bootstrap"
	appserver "servify/apps/server/internal/app/server"

	"github.com/spf13/cobra"
)

var (
	securityCheckStrict bool
)

var checkSecurityBaselineCmd = &cobra.Command{
	Use:   "check-security-baseline",
	Short: "Validate the runtime security baseline for the current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCheckSecurityBaseline(cfgFile, securityCheckStrict, cmd.OutOrStdout())
	},
	Example: strings.Join([]string{
		"  servify check-security-baseline",
		"  servify -c ../../config.production.secure.example.yml check-security-baseline --strict",
	}, "\n"),
}

func runCheckSecurityBaseline(configPath string, strict bool, out io.Writer) error {
	cfg, err := appbootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	warnings := appbootstrap.SecurityWarnings(cfg)
	router := appserver.BuildRouter(appserver.Dependencies{Config: cfg})
	warnings = append(warnings, appserver.RouteSecurityWarnings(router.Routes(), cfg)...)
	if len(warnings) == 0 {
		_, _ = fmt.Fprintln(out, "Security baseline check passed.")
		return nil
	}

	_, _ = fmt.Fprintf(out, "Security baseline check found %d issue(s):\n", len(warnings))
	for _, warning := range warnings {
		_, _ = fmt.Fprintf(out, "- %s\n", warning)
	}

	if strict {
		return fmt.Errorf("security baseline check failed in strict mode")
	}

	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Continuing because strict mode is disabled.")
	return nil
}

func init() {
	rootCmd.AddCommand(checkSecurityBaselineCmd)
	checkSecurityBaselineCmd.Flags().BoolVar(&securityCheckStrict, "strict", false, "exit non-zero when baseline warnings are present")
}
