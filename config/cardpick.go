// config/cardpick.go

package config

import (
	"errors"
	"path"
	"strconv"
	"strings"

	"github.com/vinegarhq/vinegar/sysinfo"
)

var (
	ErrOpenGLBlind = errors.New("OpenGL is not capable of choosing the right GPU, it must be explicitly defined")
	ErrNoCardFound = errors.New("GPU not found")
	ErrBadGpuIndex = errors.New("GPU index cannot be negative")
)

func (b *Binary) pickCard() error {
	if b.ForcedGpu == "" {
		return nil
	}

	n := len(sysinfo.Cards)
	idx := -1
	prime := false
	aliases := map[string]int{
		"integrated":     0,
		"prime-discrete": 1,
	}

	if i, ok := aliases[b.ForcedGpu]; ok {
		idx = i
		prime = true
	} else {
		i, err := strconv.Atoi(b.ForcedGpu)
		if err != nil {
			return err
		}

		idx = i
	}

	if prime {
		vk := b.Dxvk || b.Renderer == "Vulkan"

		if n <= 1 {
			return nil
		}

		if n > 2 && !vk {
			return ErrOpenGLBlind
		}

		if !sysinfo.Cards[0].Embedded {
			return nil
		}
	}

	if idx < 0 {
		return ErrBadGpuIndex
	}

	if n < idx+1 {
		return ErrNoCardFound
	}

	c := sysinfo.Cards[idx]

	b.Env.Set("MESA_VK_DEVICE_SELECT_FORCE_DEFAULT_DEVICE", "1")
	b.Env.Set("DRI_PRIME",
		"pci-"+strings.NewReplacer(":", "_", ".", "_").Replace(path.Base(c.Device)),
	)

	// Adjustments for FreeBSD
	// Note: Adjust these conditions based on your specific requirements for FreeBSD
	if strings.Contains(c.Driver, "nvidia") { // Workaround for OpenGL on Nvidia GPUs
		b.Env.Set("__GLX_VENDOR_LIBRARY_NAME", "nvidia")
	} else {
		b.Env.Set("__GLX_VENDOR_LIBRARY_NAME", "mesa")
	}

	return nil
}
