//go:build !windows
// +build !windows

package popup

import (
	"image/color"
	"log"
	"math/rand"
	"os/exec"
	"sultengutt/assets"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomTheme extends the default theme with larger text sizes
type CustomTheme struct {
	fyne.Theme
}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Run displays a modern, minimalist popup reminder
func RunMacPopup(orderUrl string) {
	myApp := app.New()
	myApp.Settings().SetTheme(&CustomTheme{Theme: theme.DefaultTheme()})
	window := myApp.NewWindow("Sultengutt")

	// Clean, modern window size
	window.Resize(fyne.NewSize(460, 400))

	window.CenterOnScreen()

	messages := []string{
		"Time to order surprise dinner!",
		"Save money!!!",
	}

	// Load mantra from config file
	var mantraText string
	mantraLoader, err := assets.NewMantraLoader()
	if err != nil {
		mantraText = "Stay focused and keep moving forward"
	} else {
		mantraText = mantraLoader.GetMantra()
	}

	// Large emoji
	emojiText := canvas.NewText("🍕", color.White)
	emojiText.TextSize = 48
	emojiText.Alignment = fyne.TextAlignCenter
	emojiContainer := container.NewCenter(emojiText)

	// Clean, modern title
	titleText := canvas.NewText("Surprise Dinner Reminder", color.White)
	titleText.TextSize = 24
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter
	titleContainer := container.NewCenter(titleText)

	// Simple subtitle
	subtitleText := canvas.NewText(messages[rand.Intn(len(messages))], color.RGBA{200, 200, 200, 255})
	subtitleText.TextSize = 16
	subtitleText.Alignment = fyne.TextAlignCenter
	subtitleContainer := container.NewCenter(subtitleText)

	// Mantra section header
	mantraHeaderText := canvas.NewText("Your mantra for today", color.RGBA{180, 180, 180, 255})
	mantraHeaderText.TextSize = 14
	mantraHeaderText.Alignment = fyne.TextAlignCenter
	mantraHeaderContainer := container.NewCenter(mantraHeaderText)

	// Mantra text with quotes
	mantraQuoteText := canvas.NewText("\""+mantraText+"\"", color.White)
	mantraQuoteText.TextSize = 18
	mantraQuoteText.TextStyle = fyne.TextStyle{Italic: true}
	mantraQuoteText.Alignment = fyne.TextAlignCenter

	// Wrap mantra in a container with max width for better line breaks
	mantraContainer := container.NewCenter(
		container.NewPadded(mantraQuoteText),
	)

	// Clean button styling with better sizing
	orderButton := widget.NewButton("Order Now", func() {
		cmd := exec.Command("open", orderUrl)
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Failed to open URL: %v", err)
		}
		window.Close()
	})
	orderButton.Importance = widget.HighImportance

	skipButton := widget.NewButton("Skip Today", func() {
		window.Close()
	})

	// Button container with equal spacing
	buttonContainer := container.New(
		layout.NewGridLayoutWithColumns(3),
		skipButton,
		orderButton,
	)

	// Clean vertical layout with proper spacing
	content := container.NewVBox(
		container.NewPadded(emojiContainer),
		titleContainer,
		container.NewPadded(subtitleContainer),
		container.NewPadded(widget.NewSeparator()),
		mantraHeaderContainer,
		container.NewPadded(mantraContainer),
		layout.NewSpacer(),
		container.NewPadded(buttonContainer),
	)

	// Add padding around the entire content
	paddedContent := container.NewPadded(content)

	window.SetContent(paddedContent)

	// Auto-close after 3 minutes
	go func() {
		time.Sleep(3 * time.Minute)
		if window != nil {
			window.Close()
		}
	}()

	window.ShowAndRun()
}
