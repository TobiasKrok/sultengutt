package installer

import (
	"fmt"
	"regexp"
	"strings"
	"sultengutt/internal/config"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	asciiStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
	nameStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true).MarginTop(0).Align(lipgloss.Center)
	textStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#B6E3FF")).Margin(1, 0, 0, 0)
)

const sultenguttFace = `
  ______
 /      \
|  O  O  |
|   ^    |
|  ---   |
 \______/
`

func WelcomeNote() string {
	ascii := asciiStyle.Render(sultenguttFace)
	name := nameStyle.Render("Sultengutt")
	welcome := textStyle.Render("Welcome to Sultengutt!\nNever miss a surprise dinner again. 🍔")
	return fmt.Sprintf("%s\n%s\n\n%s", ascii, name, welcome)
}

func RunInstaller(alreadyInstalled bool, prev config.InstallOptions) (config.InstallOptions, error) {
	var options config.InstallOptions = prev

	// If already installed, show config and ask if user wants to reinstall
	if alreadyInstalled {
		confString := fmt.Sprintf(
			"Your current Sultengutt config:\n\n  Days: %s\n  Hour: %s\n",
			strings.Join(prev.Days, ", "), prev.Hour,
		)
		var reinstall bool
		// Confirmation form
		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Already Installed").
					Description(
						textStyle.Render("Sultengutt is already installed!\n\n"+confString+"\nDo you want to reinstall and change your configuration?"),
					),
				huh.NewConfirm().
					Title("Reinstall?").
					Affirmative("Yes, change my config").
					Negative("No, keep as is").
					Value(&reinstall),
			),
		).WithTheme(huh.ThemeDracula())

		// Run confirmation
		if err := confirmForm.Run(); err != nil {
			return prev, err // Cancelled
		}
		if !reinstall {
			return prev, fmt.Errorf("Installation cancelled by user")
		}
	}

	// The actual main form, with prev values as defaults if re-installing
	mainForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Sultengutt Installer").
				Description(WelcomeNote()).
				Next(true).
				NextLabel("Continue"),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose Your Schedule").
				Description("Pick the days you want your dinner reminder.").
				Options(
					huh.NewOption("Monday", "Monday"),
					huh.NewOption("Tuesday", "Tuesday"),
					huh.NewOption("Wednesday", "Wednesday"),
					huh.NewOption("Thursday", "Thursday"),
					huh.NewOption("Friday", "Friday"),
					huh.NewOption("Saturday", "Saturday"),
					huh.NewOption("Sunday", "Sunday"),
				).
				Value(&options.Days).
				Validate(func(t []string) error {
					if len(t) == 0 {
						return fmt.Errorf("You must select at least one day.")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Reminder Time").
				Description("When should we remind you?\n(24-hour format, e.g. 09:30 or 17:30)").
				Value(&options.Hour).
				Validate(func(t string) error {
					time24hRegex := regexp.MustCompile(config.Time24hRegex)
					if !time24hRegex.MatchString(t) {
						return fmt.Errorf("Time must be in format HH:MM (24h)")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeDracula())

	err := mainForm.Run()
	return options, err
}
