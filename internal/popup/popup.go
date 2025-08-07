package popup

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Run displays a simple cross-platform Fyne popup reminder
// This version avoids complex OpenGL features that may not work on all Windows systems
func Run() {
	myApp := app.New()
	window := myApp.NewWindow("🍽️ Sultengutt Reminder")
	
	// Set a reasonable size and center the window
	window.Resize(fyne.NewSize(400, 300))
	window.CenterOnScreen()
	
	// Random messages
	messages := []string{
		"Your colleagues are counting on you!",
		"Time to brighten someone's day!",
		"Spread the joy with a surprise meal!",
		"Make today special for your team!",
	}
	
	tips := []string{
		"💡 Tip: Consider dietary restrictions when ordering",
		"💡 Fact: Surprise dinners boost team morale by 73%!",
		"💡 Pro tip: Pizza is always a crowd favorite",
		"💡 Remember: Variety is the spice of life!",
	}
	
	// Create simple UI elements without custom rendering
	title := widget.NewLabel("🍽️ Time for Surprise Dinner!")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	
	message := widget.NewLabel(messages[rand.Intn(len(messages))])
	message.Alignment = fyne.TextAlignCenter
	message.Wrapping = fyne.TextWrapWord
	
	tip := widget.NewLabel(tips[rand.Intn(len(tips))])
	tip.Alignment = fyne.TextAlignCenter
	tip.Wrapping = fyne.TextWrapWord
	
	// Create buttons
	orderButton := widget.NewButton("Order Now! 🛍️", func() {
		fmt.Println("Opening food ordering site...")
		window.Close()
	})
	orderButton.Importance = widget.HighImportance
	
	snoozeButton := widget.NewButton("Snooze (30 min)", func() {
		fmt.Println("Snoozed for 30 minutes")
		window.Close()
	})
	
	skipButton := widget.NewButton("Skip Today", func() {
		fmt.Println("Skipped today's reminder")
		window.Close()
	})
	
	// Create simple layout without custom widgets or animations
	content := container.NewVBox(
		container.NewPadded(title),
		widget.NewSeparator(),
		container.NewPadded(message),
		container.NewPadded(tip),
		layout.NewSpacer(),
		widget.NewSeparator(),
		container.NewGridWithColumns(3,
			skipButton,
			snoozeButton,
			orderButton,
		),
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
