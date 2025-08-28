package popup

// ShowPopup is the platform-specific popup implementation
// The actual implementation is in popup_darwin.go and popup_windows.go
func ShowPopup(siteLink string) {
	showPopup(siteLink)
}