// sysinfo/sysinfo.go
package sysinfo

import (
	"os"
	"os/exec"
	"strings"
)

var (
	Kernel    string
	CPU       Processor
	Cards     []Card
	Distro    string
	InJail    bool
	InFlatpak bool
)

func init() {
	Kernel = getKernel()
	CPU = getCPU()
	Cards = getCards()
	Distro = getDistro()

	_, err := os.Stat("/.jail")
	InJail = err == nil

	// On BSD, Flatpak is not present, so set InFlatpak to false
	InFlatpak = false
}

func getKernel() string {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

func getCPU() Processor {
	model := "Not Available" // Any ideas to implement a logic add it here

	return Processor{
		Name: model,
	}
}

func getCards() []Card {
	return []Card{}
}

func getDistro() string {
	return ""
}

type Processor struct {
	Name string
}

// Card represents information about a hardware card
type Card struct {
	Model    string
	Embedded bool   // Added field to represent whether the card is embedded
	Driver   string // Added field to represent the driver associated with the card
	Device   string // Added field to represent the device associated with the card
	Path     string

}
