package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vinegarhq/vinegar/sysinfo"
)

type Data struct {
	Project string `json:"project"`
	Distro  string `json:"distro"`
	Kernel  string `json:"kernel"`
	Flatpak bool   `json:"flatpak"`
	CPU     string `json:"cpu"`
	GPUs    string `json:"gpu"` // gpu driver list seperated by commas
}

func SubmitMerlin() error {
	var cs []string
	for _, c := range sysinfo.Cards {
		cs = append(cs, c.Driver)
	}

	d := Data{
		Project: "Vinegar " + Version,
		Distro:  sysinfo.Distro,
		Kernel:  sysinfo.Kernel,
		Flatpak: sysinfo.InFlatpak,
		CPU:     sysinfo.CPU.Name,
		GPUs:    strings.Join(cs, ","),
	}

	sd, err := json.Marshal(d)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://merlin.vinegarhq.org", "application/json", bytes.NewBuffer(sd))
	if err != nil {
		return fmt.Errorf("merlin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("merlin: bad status: %s", resp.Status)
	}

	log.Println("Successfully sent hardware information to merlin")

	return nil
}
