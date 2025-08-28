//go:build windows
// +build windows

package windows

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sultengutt/assets"
)

// Run displays the Windows popup by executing the PowerShell script
func RunWindowsPopup() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".sultengutt", "popup.ps1")

	// Execute the PowerShell script
	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-WindowStyle", "Hidden", "-File", configDir)
	cmd.Run()
}

// GenerateWindowsScript creates a modern PowerShell script for the Windows popup
func GenerateWindowsScript(configDir string, siteLink string) (string, error) {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	scriptPath := filepath.Join(configDir, "popup.ps1")

	// Random motivational messages
	messages := []string{
		"Time to order surprise dinner!",
		"Save money!!",
	}

	// Load mantra from config file
	var mantraText string
	mantraLoader, err := assets.NewMantraLoader()
	if err != nil {
		mantraText = "Stay focused and keep moving forward"
	} else {
		mantraText = mantraLoader.GetMantra()
	}

	randomMessage := messages[rand.Intn(len(messages))]

	// Modern PowerShell script with WPF for better UI
	scriptContent := fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework
Add-Type -AssemblyName System.Drawing
Add-Type -AssemblyName System.Windows.Forms

[xml]$xaml = @"
<Window
    xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
    xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
    Title="Sultengutt"
    Height="420"
    Width="480"
    WindowStartupLocation="CenterScreen"
    ResizeMode="NoResize"
    WindowStyle="SingleBorderWindow"
    Background="#FF2D2D30">
    
    <Grid Margin="20">
        <Grid.RowDefinitions>
            <RowDefinition Height="Auto"/>
            <RowDefinition Height="Auto"/>
            <RowDefinition Height="Auto"/>
            <RowDefinition Height="Auto"/>
            <RowDefinition Height="Auto"/>
            <RowDefinition Height="*"/>
            <RowDefinition Height="Auto"/>
        </Grid.RowDefinitions>
        
        <!-- Emoji -->
        <TextBlock Grid.Row="0" 
                   Text="🍕" 
                   FontSize="48" 
                   HorizontalAlignment="Center" 
                   Margin="0,10,0,10"
                   Foreground="White"/>
        
        <!-- Title -->
        <TextBlock Grid.Row="1" 
                   Text="Surprise Dinner Reminder" 
                   FontSize="24" 
                   FontWeight="Bold" 
                   HorizontalAlignment="Center" 
                   Foreground="White"
                   Margin="0,0,0,10"/>
        
        <!-- Subtitle -->
        <TextBlock Grid.Row="2" 
                   Text="%s" 
                   FontSize="16" 
                   HorizontalAlignment="Center" 
                   Foreground="#FFB4B4B4"
                   TextWrapping="Wrap"
                   Margin="0,0,0,20"/>
        
        <!-- Separator -->
        <Border Grid.Row="3" 
                Height="1" 
                Background="#FF505050" 
                Margin="40,10,40,20"/>
        
        <!-- Mantra Header -->
        <TextBlock Grid.Row="4" 
                   Text="Your mantra for today" 
                   FontSize="14" 
                   HorizontalAlignment="Center" 
                   Foreground="#FF969696"
                   Margin="0,0,0,10"/>
        
        <!-- Mantra Text -->
        <TextBlock Grid.Row="5" 
                   Text="&quot;%s&quot;" 
                   FontSize="18" 
                   FontStyle="Italic"
                   HorizontalAlignment="Center" 
                   VerticalAlignment="Center"
                   Foreground="White"
                   TextWrapping="Wrap"
                   TextAlignment="Center"
                   Margin="20,0,20,20"/>
        
        <!-- Buttons -->
        <Grid Grid.Row="6" Margin="0,10,0,0">
            <Grid.ColumnDefinitions>
                <ColumnDefinition Width="*"/>
                <ColumnDefinition Width="*"/>
                <ColumnDefinition Width="*"/>
            </Grid.ColumnDefinitions>
            
            <Button Name="SkipButton" 
                    Grid.Column="0" 
                    Content="Skip Today" 
                    Height="35" 
                    Margin="5,0,5,0"
                    Background="#FF505050"
                    Foreground="White"
                    BorderThickness="0"
                    FontSize="14"/>
            
            <Button Name="OrderButton" 
                    Grid.Column="2" 
                    Content="Order Now" 
                    Height="35" 
                    Margin="5,0,5,0"
                    Background="#FF007ACC"
                    Foreground="White"
                    BorderThickness="0"
                    FontSize="14"
                    FontWeight="Bold"/>
        </Grid>
    </Grid>
</Window>
"@

$reader = (New-Object System.Xml.XmlNodeReader $xaml)
$window = [Windows.Markup.XamlReader]::Load($reader)

# Get button references
$orderButton = $window.FindName("OrderButton")
$skipButton = $window.FindName("SkipButton")

# Add button click handlers
$orderButton.Add_Click({
    Write-Host "Opening food ordering site..."
    Start %s
    $window.Close()
})

$skipButton.Add_Click({
    Write-Host "Skipped today's reminder"
    $window.Close()
})

# Auto-close timer (3 minutes)
$timer = New-Object System.Windows.Threading.DispatcherTimer
$timer.Interval = [TimeSpan]::FromMinutes(3)
$timer.Add_Tick({
    $window.Close()
})
$timer.Start()

# Show the window
$window.ShowDialog() | Out-Null
`, randomMessage, mantraText, siteLink)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write script file: %w", err)
	}

	return scriptPath, nil
}

// GenerateWindowsFallbackScript creates a simpler Windows Forms script as fallback
func GenerateWindowsFallbackScript(configDir string, siteLink string) (string, error) {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	scriptPath := filepath.Join(configDir, "popup_fallback.ps1")

	// Random motivational messages
	messages := []string{
		"Time to order surprise dinner!",
		"Your team is counting on you",
		"Make someone's day special",
		"Spread joy with a meal",
	}

	// Load mantra from config file
	var mantraText string
	mantraLoader, err := assets.NewMantraLoader()
	if err != nil {
		mantraText = "Stay focused and keep moving forward"
	} else {
		mantraText = mantraLoader.GetMantra()
	}

	randomMessage := messages[rand.Intn(len(messages))]

	// Simpler Windows Forms script with modern styling
	scriptContent := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$form = New-Object System.Windows.Forms.Form
$form.Text = "Sultengutt"
$form.Size = New-Object System.Drawing.Size(480, 420)
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"
$form.MaximizeBox = $false
$form.MinimizeBox = $false
$form.BackColor = [System.Drawing.Color]::FromArgb(45, 45, 48)

# Emoji Label
$emojiLabel = New-Object System.Windows.Forms.Label
$emojiLabel.Text = "🍕"
$emojiLabel.Font = New-Object System.Drawing.Font("Segoe UI Emoji", 36)
$emojiLabel.Location = New-Object System.Drawing.Point(0, 20)
$emojiLabel.Size = New-Object System.Drawing.Size(460, 60)
$emojiLabel.TextAlign = "MiddleCenter"
$emojiLabel.ForeColor = [System.Drawing.Color]::White

# Title Label
$titleLabel = New-Object System.Windows.Forms.Label
$titleLabel.Text = "Surprise Dinner Reminder"
$titleLabel.Font = New-Object System.Drawing.Font("Segoe UI", 18, [System.Drawing.FontStyle]::Bold)
$titleLabel.Location = New-Object System.Drawing.Point(10, 80)
$titleLabel.Size = New-Object System.Drawing.Size(460, 35)
$titleLabel.TextAlign = "MiddleCenter"
$titleLabel.ForeColor = [System.Drawing.Color]::White

# Message Label
$messageLabel = New-Object System.Windows.Forms.Label
$messageLabel.Text = "%s"
$messageLabel.Font = New-Object System.Drawing.Font("Segoe UI", 12)
$messageLabel.Location = New-Object System.Drawing.Point(10, 120)
$messageLabel.Size = New-Object System.Drawing.Size(460, 30)
$messageLabel.TextAlign = "MiddleCenter"
$messageLabel.ForeColor = [System.Drawing.Color]::FromArgb(200, 200, 200)

# Separator
$separator = New-Object System.Windows.Forms.Label
$separator.Text = ""
$separator.Location = New-Object System.Drawing.Point(40, 160)
$separator.Size = New-Object System.Drawing.Size(400, 2)
$separator.BorderStyle = "Fixed3D"

# Mantra Header
$mantraHeader = New-Object System.Windows.Forms.Label
$mantraHeader.Text = "Your mantra for today"
$mantraHeader.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$mantraHeader.Location = New-Object System.Drawing.Point(10, 180)
$mantraHeader.Size = New-Object System.Drawing.Size(460, 25)
$mantraHeader.TextAlign = "MiddleCenter"
$mantraHeader.ForeColor = [System.Drawing.Color]::FromArgb(180, 180, 180)

# Mantra Label
$mantraLabel = New-Object System.Windows.Forms.Label
$mantraLabel.Text = '"%s"'
$mantraLabel.Font = New-Object System.Drawing.Font("Segoe UI", 13, [System.Drawing.FontStyle]::Italic)
$mantraLabel.Location = New-Object System.Drawing.Point(30, 210)
$mantraLabel.Size = New-Object System.Drawing.Size(420, 80)
$mantraLabel.TextAlign = "MiddleCenter"
$mantraLabel.ForeColor = [System.Drawing.Color]::White

# Order Button
$orderButton = New-Object System.Windows.Forms.Button
$orderButton.Text = "Order Now"
$orderButton.Font = New-Object System.Drawing.Font("Segoe UI", 10, [System.Drawing.FontStyle]::Bold)
$orderButton.Location = New-Object System.Drawing.Point(320, 320)
$orderButton.Size = New-Object System.Drawing.Size(120, 40)
$orderButton.BackColor = [System.Drawing.Color]::FromArgb(0, 122, 204)
$orderButton.ForeColor = [System.Drawing.Color]::White
$orderButton.FlatStyle = "Flat"
$orderButton.FlatAppearance.BorderSize = 0
$orderButton.Add_Click({
    Write-Host "Opening food ordering site..."
    $form.Close()
})

# Snooze Button
$snoozeButton = New-Object System.Windows.Forms.Button
$snoozeButton.Text = "Remind in 30 min"
$snoozeButton.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$snoozeButton.Location = New-Object System.Drawing.Point(180, 320)
$snoozeButton.Size = New-Object System.Drawing.Size(120, 40)
$snoozeButton.BackColor = [System.Drawing.Color]::FromArgb(112, 112, 112)
$snoozeButton.ForeColor = [System.Drawing.Color]::White
$snoozeButton.FlatStyle = "Flat"
$snoozeButton.FlatAppearance.BorderSize = 0
$snoozeButton.Add_Click({
    Write-Host "Snoozed for 30 minutes"
    $form.Close()
})

# Skip Button
$skipButton = New-Object System.Windows.Forms.Button
$skipButton.Text = "Skip Today"
$skipButton.Font = New-Object System.Drawing.Font("Segoe UI", 10)
$skipButton.Location = New-Object System.Drawing.Point(40, 320)
$skipButton.Size = New-Object System.Drawing.Size(120, 40)
$skipButton.BackColor = [System.Drawing.Color]::FromArgb(80, 80, 80)
$skipButton.ForeColor = [System.Drawing.Color]::White
$skipButton.FlatStyle = "Flat"
$skipButton.FlatAppearance.BorderSize = 0
$skipButton.Add_Click({
    Write-Host "Skipped today's reminder"
    $form.Close()
})

$form.Controls.Add($emojiLabel)
$form.Controls.Add($titleLabel)
$form.Controls.Add($messageLabel)
$form.Controls.Add($separator)
$form.Controls.Add($mantraHeader)
$form.Controls.Add($mantraLabel)
$form.Controls.Add($orderButton)
$form.Controls.Add($snoozeButton)
$form.Controls.Add($skipButton)

# Auto-close timer (3 minutes)
$timer = New-Object System.Windows.Forms.Timer
$timer.Interval = 180000
$timer.Add_Tick({
    $form.Close()
})
$timer.Start()

$form.ShowDialog()
`, randomMessage, mantraText)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write fallback script file: %w", err)
	}

	return scriptPath, nil
}
