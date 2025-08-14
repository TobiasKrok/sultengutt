package scheduler

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sultengutt/internal/config"
)

type WindowsScheduler struct {
	installOptions    config.InstallOptions
	execPath          string
	schedulerExecPath string
	configDir         string
}

func (w *WindowsScheduler) RegisterTask() error {
	_, err := w.createScriptFile()
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	
	args := w.createTask()
	cmd := exec.Command(w.schedulerExecPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to register task: %w\n%s", err, out)
	}
	return nil
}

func (w *WindowsScheduler) UnregisterTask() error {
	exists, err := w.TaskExists()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	cmd := exec.Command(w.schedulerExecPath, "/delete", "/tn", "Sultengutt", "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unregister task: %w\n%s", err, out)
	}
	return nil
}

func (w *WindowsScheduler) TaskExists() (bool, error) {
	// we assume there is only one Sultengutt task :)
	cmd := exec.Command(w.schedulerExecPath, "/query", "/tn", "Sultengutt")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (w *WindowsScheduler) createTask() []string {
	var days []string
	for _, day := range w.installOptions.Days {
		days = append(days, strings.ToUpper(day)[0:3]) // schtask accepts strings like MON,TUE,THU
	}
	scriptPath := filepath.Join(w.configDir, "popup.ps1")
	return []string{"/create",
		"/tn", "Sultengutt",
		"/tr", fmt.Sprintf("powershell.exe -ExecutionPolicy Bypass -WindowStyle Hidden -File \"%s\"", scriptPath),
		"/sc", "weekly",
		"/d", strings.Join(days, ","),
		"/st", w.installOptions.Hour,
		"/f"}
}

func (w *WindowsScheduler) createScriptFile() (string, error) {
	if err := os.MkdirAll(w.configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	
	scriptPath := filepath.Join(w.configDir, "popup.ps1")
	
	messages := []string{
		"Your colleagues are counting on you!",
		"Time to brighten someone's day!",
		"Spread the joy with a surprise meal!",
		"Make today special for your team!",
	}
	
	tips := []string{
		"Tip: Consider dietary restrictions when ordering",
		"Fact: Surprise dinners boost team morale by 73%!",
		"Pro tip: Pizza is always a crowd favorite",
		"Remember: Variety is the spice of life!",
	}
	
	randomMessage := messages[rand.Intn(len(messages))]
	randomTip := tips[rand.Intn(len(tips))]
	
	scriptContent := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$form = New-Object System.Windows.Forms.Form
$form.Text = "Sultengutt Reminder"
$form.Size = New-Object System.Drawing.Size(450, 350)
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"
$form.MaximizeBox = $false
$form.MinimizeBox = $false

$titleLabel = New-Object System.Windows.Forms.Label
$titleLabel.Text = "Time for Surprise Dinner!"
$titleLabel.Font = New-Object System.Drawing.Font("Arial", 16, [System.Drawing.FontStyle]::Bold)
$titleLabel.Location = New-Object System.Drawing.Point(10, 20)
$titleLabel.Size = New-Object System.Drawing.Size(420, 30)
$titleLabel.TextAlign = "MiddleCenter"

$messageLabel = New-Object System.Windows.Forms.Label
$messageLabel.Text = "%s"
$messageLabel.Font = New-Object System.Drawing.Font("Arial", 11)
$messageLabel.Location = New-Object System.Drawing.Point(10, 70)
$messageLabel.Size = New-Object System.Drawing.Size(420, 40)
$messageLabel.TextAlign = "MiddleCenter"

$tipLabel = New-Object System.Windows.Forms.Label
$tipLabel.Text = "%s"
$tipLabel.Font = New-Object System.Drawing.Font("Arial", 9)
$tipLabel.Location = New-Object System.Drawing.Point(10, 130)
$tipLabel.Size = New-Object System.Drawing.Size(420, 40)
$tipLabel.TextAlign = "MiddleCenter"
$tipLabel.ForeColor = [System.Drawing.Color]::DarkGray

$orderButton = New-Object System.Windows.Forms.Button
$orderButton.Text = "Order Now!"
$orderButton.Font = New-Object System.Drawing.Font("Arial", 10, [System.Drawing.FontStyle]::Bold)
$orderButton.Location = New-Object System.Drawing.Point(175, 200)
$orderButton.Size = New-Object System.Drawing.Size(100, 40)
$orderButton.BackColor = [System.Drawing.Color]::LightGreen
$orderButton.Add_Click({
    Write-Host "Opening food ordering site..."
    $form.Close()
})

$snoozeButton = New-Object System.Windows.Forms.Button
$snoozeButton.Text = "Snooze 30m"
$snoozeButton.Location = New-Object System.Drawing.Point(50, 200)
$snoozeButton.Size = New-Object System.Drawing.Size(100, 40)
$snoozeButton.Add_Click({
    Write-Host "Snoozed for 30 minutes"
    $form.Close()
})

$skipButton = New-Object System.Windows.Forms.Button
$skipButton.Text = "Skip Today"
$skipButton.Location = New-Object System.Drawing.Point(300, 200)
$skipButton.Size = New-Object System.Drawing.Size(100, 40)
$skipButton.Add_Click({
    Write-Host "Skipped today's reminder"
    $form.Close()
})

$form.Controls.Add($titleLabel)
$form.Controls.Add($messageLabel)
$form.Controls.Add($tipLabel)
$form.Controls.Add($orderButton)
$form.Controls.Add($snoozeButton)
$form.Controls.Add($skipButton)

$timer = New-Object System.Windows.Forms.Timer
$timer.Interval = 180000
$timer.Add_Tick({
    $form.Close()
})
$timer.Start()

$form.ShowDialog()
`, randomMessage, randomTip)
	
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write script file: %w", err)
	}
	
	return scriptPath, nil
}
