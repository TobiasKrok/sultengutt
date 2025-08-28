//go:build darwin
// +build darwin

package popup

import (
	macpop "sultengutt/internal/popup/mac"
)

func showPopup(siteLink string) {
	macpop.RunMacPopup(siteLink)
}
