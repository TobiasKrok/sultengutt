package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sultengutt/internal/config"
	"sultengutt/internal/installer"
	"sultengutt/internal/popup"
	"sultengutt/internal/scheduler"
	"sultengutt/internal/utils"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Styles for consistent formatting
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#27AE60"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E74C3C"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3498DB"))
)

func main() {

	cm, err := config.NewConfigManager()
	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration: %w", err))
	}
	cfg, err := cm.Load()

	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	rootCmd := &cobra.Command{
		Use:   "sultengutt",
		Short: "Never miss a surprise dinner again.",
		Long: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("Sultengutt") + "\n\n" +
			"Be reminded to order your surprise dinner on your schedule! Never be hungry again 🍔",
		Example: `  sultengutt install
  sultengutt execute
  sultengutt pause 1 day
  sultengutt resume
  sultengutt status`,
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up Sultengutt with interactive installer",
		Long:  "Install or reinstall Sultengutt with an interactive installer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(cfg, cm)
		},
	}

	executeCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute Sultengutt reminder",
		Long:  "Executes Sultengutt to trigger the popup reminder.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.IsPaused() && cfg.PausedUntil == 0 || cfg.IsPaused() && cfg.PausedUntil > 0 && time.Now().Unix() < cfg.PausedUntil {
				fmt.Println("paused. Use 'sultengutt resume' to unpause.")
				return nil
			}
			// check if we need to resume
			if cfg.IsPaused() && cfg.PausedUntil > 0 && time.Now().Unix() >= cfg.PausedUntil {
				cfg.PausedUntil = -1
				err := cm.Save(cfg)
				if err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
			}

			popup.ShowPopup(cfg.InstallOptions.SiteLink)

			return nil
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the Sultengutt reminder. Use -help for full examples",
		Long:  "Pause Sultengutt for a period of time or until you resume\n\nAllowed units: day(s), week(s), month(s) or indefinitely (no arguments) ",
		Example: `sultengutt pause 1 day
	sultengutt pause 4 weeks
	sultengutt pause 1 month
	sultengutt pause // Pause indefinitely`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.IsFreshInstall() {
				fmt.Println(errorStyle.Render("Sultengutt has not been installed yet, please run `sultengutt install` first"))
				return nil
			}
			if err := runPause(args, cfg); err != nil {
				return err
			}
			return cm.Save(cfg)
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume Sultengutt reminders",
		Long:  "Manually resume Sultengutt reminders.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.PausedUntil = -1
			err := cm.Save(cfg)
			if err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println(infoStyle.Render("Resumed Sultengutt reminders"))
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show current status of Sultengutt",
		Long:  "Show current status of Sultengutt.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.IsFreshInstall() {
				fmt.Println(errorStyle.Render("Sultengutt has not been installed yet, please run `sultengutt install` first"))
				return
			}
			runStatus(*cfg)
		},
	}

	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Sultengutt completely",
		Long:  "Uninstall Sultengutt, removing all scheduled tasks and configuration files.",
		Example: `  sultengutt uninstall
  sultengutt uninstall --confirm`,
		RunE: func(cmd *cobra.Command, args []string) error {
			confirm, _ := cmd.Flags().GetBool("confirm")
			return runUninstall(cfg, cm, confirm)
		},
	}
	uninstallCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")

	rootCmd.AddCommand(installCmd, executeCmd, pauseCmd, resumeCmd, statusCmd, uninstallCmd)

	rootCmd.SetErrPrefix(errorStyle.Render("Error:"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorStyle.Render("✗ " + err.Error()))
		os.Exit(1)
	}
}
func runInstall(cfg *config.Config, cm *config.ConfigManager) error {

	installed := cfg.IsFreshInstall()
	opts, err := installer.RunInstaller(installed, cfg.InstallOptions)
	if err != nil {
		return fmt.Errorf("installation cancelled or failed: %w", err)
	}

	cfg.InstallOptions = opts
	err = cm.Save(cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	sch := scheduler.NewScheduler(opts, cm.ConfigDir())
	if installed {
		if err := sch.UnregisterTask(); err != nil {
			return fmt.Errorf("failed to unregister old task: %w", err)
		}
	}
	if err := sch.RegisterTask(); err != nil {
		return fmt.Errorf("failed to register task: %w", err)
	}

	fmt.Println(successStyle.Render("✓ Sultengutt is set up! Happy dining!"))
	fmt.Println(infoStyle.Render("Config saved to: " + cfg.Path() + "\n\n"))
	return nil
}

func runPause(args []string, cfg *config.Config) error {
	if len(args) == 0 {
		// Pause indefinitely
		cfg.PausedUntil = 0
		fmt.Println("Paused indefinitely. Use 'sultengutt resume' to unpause.")
		return nil
	}

	duration, err := utils.ParseDuration(args)
	if err != nil {
		return fmt.Errorf("error parsing duration: %v", err)
	}

	scheduledHour := cfg.InstallOptions.Hour
	pauseUntil, err := utils.CalculatePauseUntil(duration, scheduledHour)
	if err != nil {
		return fmt.Errorf("error calculating pause time: %v", err)
	}

	cfg.PausedUntil = pauseUntil

	unpauseTime := time.Unix(pauseUntil, 0)
	fmt.Printf("Paused until %s at %s\n",
		unpauseTime.Format("Monday, January 2, 2006"),
		unpauseTime.Format("15:04"))

	return nil
}

func runStatus(cfg config.Config) {
	fmt.Println("╭─────────────────────────────────────╮")
	fmt.Println("│      SULTENGUTT STATUS              │")
	fmt.Println("╰─────────────────────────────────────╯")
	fmt.Println()
	fmt.Println("Schedule:")
	fmt.Println("  Hour: " + cfg.InstallOptions.Hour)
	fmt.Println("  Days: " + strings.Join(cfg.InstallOptions.Days, ", "))
	fmt.Println()
	fmt.Println("Status:")
	fmt.Println("  Config path: " + cfg.Path())
	if cfg.PausedUntil > 0 {
		fmt.Println("  Paused: paused until " + time.Unix(cfg.PausedUntil, 0).Format("Monday, January 2, 2006 15:04"))
		fmt.Println("  tip: use 'sultengutt resume' to unpause early")
	} else {
		fmt.Println("  Paused: not paused (active)")
	}

}

func runUninstall(cfg *config.Config, cm *config.ConfigManager, skipConfirm bool) error {
	if cfg.IsFreshInstall() {
		fmt.Println(infoStyle.Render("Sultengutt is not installed"))
		return nil
	}

	if !skipConfirm {
		fmt.Println("This will completely uninstall Sultengutt.")
		fmt.Println(errorStyle.Render("WARNING: This will remove all scheduled tasks and configuration files."))
		fmt.Print("\nAre you sure you want to continue? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Uninstall cancelled")
			return nil
		}
	}
	sch := scheduler.NewScheduler(cfg.InstallOptions, cm.ConfigDir())
	if err := sch.UnregisterTask(); err != nil {
		return fmt.Errorf("failed to unregister scheduled task: %w", err)
	}
	fmt.Println(successStyle.Render("Removed scheduled tasks"))

	if err := cm.Clean(); err != nil {
		return fmt.Errorf("failed to clean configuration: %w", err)
	}
	fmt.Println(successStyle.Render("Removed configuration files"))

	fmt.Println(successStyle.Render("\n✓ Sultengutt has been uninstalled"))

	return nil
}
