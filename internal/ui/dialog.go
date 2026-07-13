package ui

import (
	"os/exec"
	"strings"
)

// OpenFileDialog uses PowerShell to open a native Windows file dialog.
func OpenFileDialog(filter string) (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.OpenFileDialog
$f.Filter = "` + filter + `"
$f.ShowHelp = $false
$res = $f.ShowDialog()
if ($res -eq "OK") { Write-Output $f.FileName }
`
	out, err := exec.Command("powershell", "-Sta", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// OpenFolderDialog uses PowerShell to open a native Windows folder dialog.
func OpenFolderDialog() (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$res = $f.ShowDialog()
if ($res -eq "OK") { Write-Output $f.SelectedPath }
`
	out, err := exec.Command("powershell", "-Sta", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
