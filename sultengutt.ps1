# SULTENGUTT

Import-Module "$PSScriptRoot\sultengutt.psm1" -Force

Add-Type -AssemblyName System.Windows.Forms


# Sultengutt ASCII art
$cmd = @"
       .-""""""-.
     .'          '.
    /   O      O   \
   :           `    :
   |                |   
   :    .------.    :
    \  '        '  /
     '.          .'
       '-......-'

       SULTENGUTT


"@

Write-Host $cmd
Write-Host @"
Welcome to Sultengutt!

This script will set up a daily task to remind you to order Suprise Dinner.
The task will run at 9:30 AM every day and will display a message box to remind you.
"@ -ForegroundColor Cyan


$checkinPath = "$env:USERPROFILE\.sultengutt"
New-SultenguttSetup

Copy-Item -Path "$PSScriptRoot\sultengutt.psm1" -Destination "$checkinPath\tasks\sultengutt.psm1" -Force

#### SETUP TASKS ####

# Create daily task script content
$dailyTaskScript = @"
Import-Module "$env:USERPROFILE\.sultengutt\tasks\sultengutt.psm1" -Force
New-SultenguttSetup

`$checkinPath = "$env:USERPROFILE\.sultengutt\checkin"
Set-Content -Path `$checkinPath -Value "no" -Force

`$result = Show-SultenguttReminder

if (`$result -eq [System.Windows.Forms.DialogResult]::Yes) {
    Set-Content -Path `$checkinPath -Value "yes" -Force
    Write-Host "User has ordered dinner!" -ForegroundColor Green
} else {
    Register-SultenguttSnooze -DelayMinutes 5
    Write-Host "Snooze activated for 5 minutes" -ForegroundColor Yellow
}
"@

# Write the daily task script to a file
$dailyScriptPath = "$checkinPath\tasks\daily.ps1"
$dailyTaskScript | Out-File -FilePath $dailyScriptPath -Encoding UTF8 -Force

Write-Host "Setting up Sultengutt task..." -ForegroundColor Yellow

# Register the daily task
Register-SultenguttTask -TaskName "Sultengutt" -ScriptPath $dailyScriptPath -TriggerTime "09:30" -Description "Daily Sultengutt"

# Show the reminder dialog immediately for testing
Write-Host "`nShowing reminder dialog for testing..." -ForegroundColor Cyan
$result = Show-SultenguttReminder

if ($result -eq [System.Windows.Forms.DialogResult]::Yes) {
    $checkin = "$checkinPath\checkin"
    Set-Content -Path $checkin -Value "yes" -Force
    Write-Host "✅ Great! You've ordered your dinner!" -ForegroundColor Green
} else {
    Register-SultenguttSnooze -DelayMinutes 5
    Write-Host "⏰ Okay, I'll remind you again in 5 minutes!" -ForegroundColor Yellow
}


