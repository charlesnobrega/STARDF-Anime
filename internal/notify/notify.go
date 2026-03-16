// Package notify provides platform-agnostic desktop notifications.
// Uses only stdlib + process spawning — no CGO, no external lib required.
package notify

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Send sends a desktop notification with a title and body.
// Gracefully no-ops on unsupported platforms.
func Send(title, body string) {
	switch runtime.GOOS {
	case "windows":
		sendWindows(title, body)
	case "darwin":
		sendMacOS(title, body)
	case "linux":
		sendLinux(title, body)
	}
}

// sendWindows uses PowerShell's BurntToast-style notification
func sendWindows(title, body string) {
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom, ContentType = WindowsRuntime] | Out-Null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$template.GetElementsByTagName('text')[0].AppendChild($template.CreateTextNode('%s')) | Out-Null
$template.GetElementsByTagName('text')[1].AppendChild($template.CreateTextNode('%s')) | Out-Null
$toast = [Windows.UI.Notifications.ToastNotification]::new($template)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('StarDF-Anime').Show($toast)
`, escapePS(title), escapePS(body))

	_ = exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script).Start()
}

// sendMacOS uses osascript
func sendMacOS(title, body string) {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, body, title)
	_ = exec.Command("osascript", "-e", script).Start()
}

// sendLinux tries notify-send (libnotify)
func sendLinux(title, body string) {
	_ = exec.Command("notify-send",
		"-a", "StarDF-Anime",
		"-i", "media-playback-start",
		"--", title, body,
	).Start()
}

// escapePS escapes single quotes for PowerShell strings
func escapePS(s string) string {
	result := []rune{}
	for _, r := range s {
		if r == '\'' {
			result = append(result, '\'', '\'')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
