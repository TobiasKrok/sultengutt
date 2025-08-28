//go:build windows
// +build windows

package popup

import (
	winpop "sultengutt/internal/popup/windows"
)

func showPopup(siteLink string) {
	winpop.RunWindowsPopup()
}
