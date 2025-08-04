package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
	"sultengutt/internal/config"
	"sultengutt/internal/installer"
	"sultengutt/internal/scheduler"
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

func runInstall() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to check if Sultengutt is already installed: %w", err)
	}
	installed, err := cfg.IsFreshInstall()
	if err != nil {
		return fmt.Errorf("failed to check if Sultengutt is already installed: %w", err)
	}
	opts, err := installer.RunInstaller(installed, cfg.GetInstallOptions())
	if err != nil {
		return fmt.Errorf("installation cancelled or failed: %w", err)
	}

	err = cfg.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	sch := scheduler.NewScheduler(opts)
	if installed {
		if err := sch.UnregisterTask(); err != nil {
			return fmt.Errorf("failed to unregister old task: %w", err)
		}
	}
	if err := sch.RegisterTask(); err != nil {
		return fmt.Errorf("failed to register task: %w", err)
	}

	fmt.Println(successStyle.Render("✓ Sultengutt is set up! Happy dining!"))
	fmt.Println(infoStyle.Render("Config saved to: " + configPath + "\n\n"))
	return nil
}

func main() {

	cm, err := config.NewConfigManager()
	cfg, err := cm.Load()

	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration: %w", err))
	}

	rootCmd := &cobra.Command{
		Use:   "sultengutt",
		Short: "Never miss a surprise dinner again.",
		Long: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("Sultengutt") + "\n\n" +
			"Be reminded to order your suprise dinner on your schedule! Never be hungry again 🍔",
		Example: `  sultengutt install
  sultengutt execute
  sultengutt pause
  sultengutt resume`,
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up Sultengutt with interactive installer",
		Long:  "Install or reinstall Sultengutt with an interactive installer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall()
		},
	}

	executeCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute Sultengutt reminder",
		Long:  "Executes Sultengutt to trigger the popup reminder.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(infoStyle.Render("Execute functionality not yet implemented"))
			return nil
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the Sultengutt reminder",
		Long:  "Pause Sultengutt for a period of time or until you resume",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement pause functionality
			fmt.Println(infoStyle.Render("Pause functionality not yet implemented"))
			return nil
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume Sultengutt reminders",
		Long:  "Manually resume Sultengutt reminders.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement pause functionality
			fmt.Println(infoStyle.Render("Pause functionality not yet implemented"))
			return nil
		},
	}
	// Add commands to root
	rootCmd.AddCommand(installCmd, executeCmd, pauseCmd, resumeCmd)

	// Custom error handling with lipgloss styling
	rootCmd.SetErrPrefix(errorStyle.Render("Error:"))

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorStyle.Render("✗ " + err.Error()))
		os.Exit(1)
	}
} //a := app.New()
//w := a.NewWindow("Sultengutt")
//hello := widget.NewLabel("Remember to buy your suprise dinner for today!")
//w.SetContent(container.NewVBox(hello))
//w.ShowAndRun(
