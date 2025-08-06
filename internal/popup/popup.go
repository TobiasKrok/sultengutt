package popup

// THANK YOU CLAUDE
import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Custom theme for modern look
type sultenguttTheme struct{}

func (m sultenguttTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 18, G: 18, B: 24, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 99, G: 102, B: 241, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 243, G: 244, B: 246, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 129, G: 140, B: 248, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 79, G: 70, B: 229, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m sultenguttTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m sultenguttTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m sultenguttTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNamePadding:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Animated emoji that bounces
type bouncingEmoji struct {
	widget.BaseWidget
	emoji     string
	offset    float32
	direction float32
	animation *fyne.Animation
}

func newBouncingEmoji(emoji string) *bouncingEmoji {
	b := &bouncingEmoji{
		emoji:     emoji,
		offset:    0,
		direction: 1,
	}
	b.ExtendBaseWidget(b)

	// Create bouncing animation
	b.animation = canvas.NewPositionAnimation(
		fyne.NewPos(0, 0),
		fyne.NewPos(0, 20),
		time.Millisecond*1500,
		func(p fyne.Position) {
			b.offset = p.Y
			b.Refresh()
		},
	)
	b.animation.RepeatCount = fyne.AnimationRepeatForever
	b.animation.AutoReverse = true
	b.animation.Curve = fyne.AnimationEaseInOut

	return b
}

func (b *bouncingEmoji) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(b.emoji, color.White)
	text.TextSize = 48
	text.Alignment = fyne.TextAlignCenter

	return &bouncingEmojiRenderer{
		emoji:  text,
		widget: b,
	}
}

func (b *bouncingEmoji) Start() {
	b.animation.Start()
}

func (b *bouncingEmoji) Stop() {
	b.animation.Stop()
}

type bouncingEmojiRenderer struct {
	emoji  *canvas.Text
	widget *bouncingEmoji
}

func (r *bouncingEmojiRenderer) Layout(size fyne.Size) {
	r.emoji.Resize(size)
	r.emoji.Move(fyne.NewPos(0, r.widget.offset))
}

func (r *bouncingEmojiRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 80)
}

func (r *bouncingEmojiRenderer) Refresh() {
	r.emoji.Move(fyne.NewPos(0, r.widget.offset))
	canvas.Refresh(r.emoji)
}

func (r *bouncingEmojiRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.emoji}
}

func (r *bouncingEmojiRenderer) Destroy() {}

// Gradient background
type gradientBackground struct {
	widget.BaseWidget
}

func newGradientBackground() *gradientBackground {
	g := &gradientBackground{}
	g.ExtendBaseWidget(g)
	return g
}

func (g *gradientBackground) CreateRenderer() fyne.WidgetRenderer {
	return &gradientRenderer{widget: g}
}

type gradientRenderer struct {
	widget *gradientBackground
	raster *canvas.Raster
}

func (r *gradientRenderer) Layout(size fyne.Size) {
	if r.raster != nil {
		r.raster.Resize(size)
	}
}

func (r *gradientRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

func (r *gradientRenderer) Refresh() {
	canvas.Refresh(r.widget)
}

func (r *gradientRenderer) Objects() []fyne.CanvasObject {
	if r.raster == nil {
		r.raster = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
			// Create a gradient from purple to blue
			progress := float64(y) / float64(h)
			r := uint8(99 + progress*30)
			g := uint8(102 - progress*32)
			b := uint8(241 - progress*20)
			return color.NRGBA{R: r, G: g, B: b, A: 255}
		})
	}
	return []fyne.CanvasObject{r.raster}
}

func (r *gradientRenderer) Destroy() {}

// Particle effect for celebration
type particle struct {
	x, y   float32
	vx, vy float32
	life   float32
	emoji  string
}

type particleEffect struct {
	widget.BaseWidget
	particles []*particle
	animation *fyne.Animation
}

func newParticleEffect() *particleEffect {
	p := &particleEffect{
		particles: make([]*particle, 0),
	}
	p.ExtendBaseWidget(p)

	// Create particles
	emojis := []string{"🎉", "✨", "🎊", "⭐", "💫"}
	for i := 0; i < 15; i++ {
		p.particles = append(p.particles, &particle{
			x:     rand.Float32() * 400,
			y:     rand.Float32()*50 + 250,
			vx:    (rand.Float32() - 0.5) * 2,
			vy:    -rand.Float32()*3 - 1,
			life:  1.0,
			emoji: emojis[rand.Intn(len(emojis))],
		})
	}

	// Animate particles
	start := time.Now()
	p.animation = &fyne.Animation{
		Duration: time.Second * 3,
		Tick: func(progress float32) {
			elapsed := time.Since(start).Seconds()
			for _, particle := range p.particles {
				particle.x += particle.vx
				particle.y += particle.vy + float32(elapsed)*0.5 // gravity
				particle.life = 1.0 - progress
			}
			p.Refresh()
		},
	}

	return p
}

func (p *particleEffect) Start() {
	p.animation.Start()
}

func (p *particleEffect) CreateRenderer() fyne.WidgetRenderer {
	return &particleRenderer{widget: p}
}

type particleRenderer struct {
	widget  *particleEffect
	objects []fyne.CanvasObject
}

func (r *particleRenderer) Layout(size fyne.Size) {}

func (r *particleRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

func (r *particleRenderer) Refresh() {
	r.objects = nil
	for _, p := range r.widget.particles {
		if p.life > 0 {
			text := canvas.NewText(p.emoji, color.White)
			text.TextSize = 20
			text.Move(fyne.NewPos(p.x, p.y))
			// Fade out
			alpha := uint8(255 * p.life)
			text.Color = color.NRGBA{R: 255, G: 255, B: 255, A: alpha}
			r.objects = append(r.objects, text)
		}
	}
	canvas.Refresh(r.widget)
}

func (r *particleRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *particleRenderer) Destroy() {}

// Create the popup window
func createSultenguttPopup() {
	myApp := app.New()
	myApp.Settings().SetTheme(&sultenguttTheme{})

	window := myApp.NewWindow("Sultengutt Reminder")
	window.Resize(fyne.NewSize(450, 380))
	window.CenterOnScreen()

	// Remove window decorations for a cleaner look (optional)
	// window.SetFixedSize(true)

	// Create gradient background
	gradient := newGradientBackground()

	// Create bouncing dinner emoji
	dinnerEmoji := newBouncingEmoji("🍽️")
	dinnerEmoji.Start()

	// Create particle effects
	particles := newParticleEffect()

	// Title with animation
	title := widget.NewLabelWithStyle(
		"Time for Surprise Dinner!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	title.TextStyle.Bold = true

	// Animated subtitle
	messages := []string{
		"Your colleagues are counting on you! 🎉",
		"Time to brighten someone's day! ✨",
		"Spread the joy with a surprise meal! 🌟",
		"Make today special for your team! 💫",
	}
	subtitle := widget.NewLabel(messages[rand.Intn(len(messages))])
	subtitle.Alignment = fyne.TextAlignCenter

	// Fun facts or tips
	tips := []string{
		"💡 Tip: Consider dietary restrictions when ordering",
		"💡 Fact: Surprise dinners boost team morale by 73%!",
		"💡 Pro tip: Pizza is always a crowd favorite",
		"💡 Remember: Variety is the spice of life!",
	}
	tipLabel := widget.NewLabel(tips[rand.Intn(len(tips))])
	tipLabel.Alignment = fyne.TextAlignCenter
	tipLabel.Wrapping = fyne.TextWrapWord

	// Progress indicator (shows time until next reminder)
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0.7) // Example: 70% through the day

	progressLabel := widget.NewLabel("Next reminder: Tomorrow at 09:45")
	progressLabel.Alignment = fyne.TextAlignCenter

	// Action buttons with custom styling
	orderButton := widget.NewButton("Order Now! 🛍️", func() {
		// TODO: Open ordering website or app
		fmt.Println("Opening food ordering site...")
		particles.Start() // Celebrate the decision!

		// Close after a delay
		time.AfterFunc(time.Second*2, func() {
			window.Close()
		})
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
	skipButton.Importance = widget.LowImportance

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			container.NewCenter(dinnerEmoji),
			widget.NewSeparator(),
			title,
			subtitle,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			tipLabel,
			progressBar,
			progressLabel,
			container.NewGridWithColumns(3,
				skipButton,
				snoozeButton,
				orderButton,
			),
		),
		nil,
		nil,
		container.NewMax(
			gradient,
			particles,
		),
	)

	// Add padding
	paddedContent := container.NewPadded(content)

	window.SetContent(paddedContent)

	go func() {
		time.Sleep(3 * time.Minute)
		if window != nil {
			window.Close()
		}
	}()

	window.ShowAndRun()
}

func Run() {
	createSultenguttPopup()
}
