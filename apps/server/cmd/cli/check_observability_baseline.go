package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	appbootstrap "servify/apps/server/internal/app/bootstrap"

	"github.com/spf13/cobra"
)

var (
	observabilityCheckStrict bool
)

var checkObservabilityBaselineCmd = &cobra.Command{
	Use:   "check-observability-baseline",
	Short: "Validate the runtime observability baseline for the current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCheckObservabilityBaseline(cfgFile, observabilityCheckStrict, os.Getenv("SERVIFY_REPO_ROOT"), cmd.OutOrStdout())
	},
	Example: strings.Join([]string{
		"  servify check-observability-baseline",
		"  servify -c ../../config.production.secure.example.yml check-observability-baseline --strict",
	}, "\n"),
}

func runCheckObservabilityBaseline(configPath string, strict bool, repoRoot string, out io.Writer) error {
	cfg, err := appbootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	warnings := appbootstrap.ObservabilityWarnings(cfg, repoRoot)
	if len(warnings) == 0 {
		_, _ = fmt.Fprintln(out, "Observability baseline check passed.")
		return nil
	}

	_, _ = fmt.Fprintf(out, "Observability baseline check found %d issue(s):\n", len(warnings))
	for _, warning := range warnings {
		_, _ = fmt.Fprintf(out, "- %s\n", warning)
	}

	if strict {
		return fmt.Errorf("observability baseline check failed in strict mode")
	}

	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintln(out, "Continuing because strict mode is disabled.")
	return nil
}

func init() {
	rootCmd.AddCommand(checkObservabilityBaselineCmd)
	checkObservabilityBaselineCmd.Flags().BoolVar(&observabilityCheckStrict, "strict", false, "exit non-zero when baseline warnings are present")
}
