# TaskManager Module for Sultengutt
# Handles registration and cleanup of scheduled tasks

function Register-SultenguttTask {
    param(
        [Parameter(Mandatory=$true)]
        [string]$TaskName,
        
        [Parameter(Mandatory=$true)]
        [string]$ScriptPath,
        
        [Parameter(Mandatory=$true)]
        [string]$TriggerTime,
        
        [string]$Description = "Sultengutt task",
        [switch]$Once
    )
    
    $taskExists = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    
    if (-not $taskExists) {
        $action = New-ScheduledTaskAction -Execute "PowerShell.exe" -Argument "-WindowStyle Hidden -File `"$ScriptPath`""
        
        if ($Once) {
            $trigger = New-ScheduledTaskTrigger -Once -At (Get-Date $TriggerTime)
            $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -DeleteExpiredTaskAfter "00:01:00"
        } else {
            $trigger = New-ScheduledTaskTrigger -Daily -At $TriggerTime
            $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable
        }
        
        try {
            Register-ScheduledTask -TaskName $TaskName -Action $action -Trigger $trigger -Settings $settings -Description $Description | Out-Null
            $timeType = if ($Once) { "once at" } else { "daily at" }
            Write-Host "✅ Task '$TaskName' registered successfully ($timeType $TriggerTime)" -ForegroundColor Green
            return $true
        }
        catch {
            Write-Host "❌ Failed to register task '$TaskName': $($_.Exception.Message)" -ForegroundColor Red
            return $false
        }
    }
    else {
        Write-Host "ℹ️ Task '$TaskName' already exists" -ForegroundColor Yellow
        return $true
    }
}

function Remove-SultenguttTask {
    param(
        [Parameter(Mandatory=$true)]
        [string]$TaskName
    )
    
    $taskExists = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    
    if ($taskExists) {
        try {
            Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
            Write-Host "✅ Task '$TaskName' removed successfully" -ForegroundColor Green
            return $true
        }
        catch {
            Write-Host "❌ Failed to remove task '$TaskName': $($_.Exception.Message)" -ForegroundColor Red
            return $false
        }
    }
    else {
        Write-Host "ℹ️ Task '$TaskName' does not exist" -ForegroundColor Yellow
        return $true
    }
}

function Get-SultenguttTasks {
    param(
        [string]$TaskNamePattern = "Sultengutt*"
    )
    
    return Get-ScheduledTask -TaskName $TaskNamePattern -ErrorAction SilentlyContinue
}

function Test-SultenguttTask {
    param(
        [Parameter(Mandatory=$true)]
        [string]$TaskName
    )
    
    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    return $null -ne $task
}

function New-SultenguttSetup {
    $checkinPath = "$env:USERPROFILE\.sultengutt"
    if (-not (Test-Path $checkinPath)) {
        New-Item -ItemType Directory -Path $checkinPath | Out-Null
    }
    if (-not (Test-Path "$checkinPath\tasks")) {
        New-Item -ItemType Directory -Path "$checkinPath\tasks" | Out-Null
    }
    if (-not (Test-Path "$checkinPath\checkin")) {
        New-Item -ItemType File -Path "$checkinPath\checkin" | Out-Null
        Set-Content -Path "$checkinPath\checkin" -Value "no" -Force
    }
}

function Show-SultenguttReminder {
    Add-Type -AssemblyName System.Windows.Forms
    Add-Type -AssemblyName System.Drawing
    
    $form = New-Object System.Windows.Forms.Form
    $form.Text = "Sultengutt"
    $form.Size = New-Object System.Drawing.Size(450, 250)
    $form.StartPosition = "CenterScreen"
    $form.TopMost = $true
    $form.FormBorderStyle = "FixedDialog"
    $form.MaximizeBox = $false
    $form.MinimizeBox = $false
    
    # Create label for the message
    $label = New-Object System.Windows.Forms.Label
    $label.Location = New-Object System.Drawing.Point(20, 30)
    $label.Size = New-Object System.Drawing.Size(400, 60)
    $label.Font = New-Object System.Drawing.Font("Microsoft Sans Serif", 14, [System.Drawing.FontStyle]::Bold)
    $label.Text = "Bestill suprise dinner!!!"
    $label.TextAlign = "MiddleCenter"
    $label.ForeColor = [System.Drawing.Color]::DarkRed
    $form.Controls.Add($label)
    
    # Create "I've ordered" button
    $orderedButton = New-Object System.Windows.Forms.Button
    $orderedButton.Location = New-Object System.Drawing.Point(50, 130)
    $orderedButton.Size = New-Object System.Drawing.Size(150, 40)
    $orderedButton.Text = "I've ordered"
    $orderedButton.Font = New-Object System.Drawing.Font("Microsoft Sans Serif", 10)
    $orderedButton.BackColor = [System.Drawing.Color]::LightGreen
    $orderedButton.DialogResult = [System.Windows.Forms.DialogResult]::Yes
    $form.Controls.Add($orderedButton)
    
    # Create "Remind me in 5 min" button
    $snoozeButton = New-Object System.Windows.Forms.Button
    $snoozeButton.Location = New-Object System.Drawing.Point(250, 130)
    $snoozeButton.Size = New-Object System.Drawing.Size(150, 40)
    $snoozeButton.Text = "Remind me 5 min"
    $snoozeButton.Font = New-Object System.Drawing.Font("Microsoft Sans Serif", 10)
    $snoozeButton.BackColor = [System.Drawing.Color]::LightYellow
    $snoozeButton.DialogResult = [System.Windows.Forms.DialogResult]::No
    $form.Controls.Add($snoozeButton)
    
    # Show the form and get result
    $result = $form.ShowDialog()
    $form.Dispose()
    
    return $result
}

function Register-SultenguttSnooze {
    param(
        [int]$DelayMinutes = 5
    )
    $snoozeScript = @"
# Import module from the original source location
Import-Module "$env:USERPROFILE\.sultengutt\tasks\sultengutt.psm1" -Force
`$checkinPath = "$env:USERPROFILE\.sultengutt\checkin"
`$checkin = Get-Content `$checkinPath -ErrorAction SilentlyContinue

if (`$checkin -eq "no") {
    `$result = Show-SultenguttReminder
    
    if (`$result -eq [System.Windows.Forms.DialogResult]::Yes) {
        Set-Content -Path `$checkinPath -Value "yes" -Force
        Remove-SultenguttTask -TaskName "SultenguttSnooze"
    } else {
        Register-SultenguttSnooze -DelayMinutes 5
    }
} else {
    Remove-SultenguttTask -TaskName "SultenguttSnooze"
}
"@
    
    # Write snooze script and register task
    $snoozeScriptPath = "$env:TEMP\SultenguttSnooze.ps1"
    $snoozeScript | Out-File -FilePath $snoozeScriptPath -Encoding UTF8 -Force
    
    # Remove existing snooze task if it exists
    Remove-SultenguttTask -TaskName "SultenguttSnooze"

    $snoozeTime = (Get-Date).AddMinutes($DelayMinutes)
    return Register-SultenguttTask -TaskName "SultenguttSnooze" -ScriptPath $snoozeScriptPath -TriggerTime $snoozeTime.ToString("HH:mm") -Description "Sultengutt snooze reminder" -Once
}

Export-ModuleMember -Function Register-SultenguttTask, Remove-SultenguttTask, Get-SultenguttTasks, Test-SultenguttTask, New-SultenguttSetup, Show-SultenguttReminder, Register-SultenguttSnooze