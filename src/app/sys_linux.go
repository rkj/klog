//go:build linux

package app

var POTENTIAL_EDITORS = []string{"vim", "vi", "nano", "pico"}

var POTENTIAL_FILE_EXLORERS = []string{"xdg-open"}

func flagAsHidden(path string) {
	// Nothing to do on UNIX due to the dotfile convention
}