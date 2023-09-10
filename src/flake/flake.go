package flake

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/hymkor/go-lazy"
)

var (
	BrandedRegistry = "__std" // keep for now for historic reasons
)

type FlakeRegistry struct {
	FlakeRef      string
	Registry      string
	FlakeRegistry string
}

func LocalFlakeRegistry() FlakeRegistry {
	return FlakeRegistry{
		FlakeRef:      "./",
		Registry:      BrandedRegistry,
		FlakeRegistry: fmt.Sprintf("%s#%s", "./", BrandedRegistry),
	}
}

func (r *FlakeRegistry) InitFlakeRef(system string) string {
	return fmt.Sprintf("%s.init.%s", r.FlakeRegistry, system)
}

func (r *FlakeRegistry) RefCellsFrom() string {
	return fmt.Sprintf("%s.cellsFrom", r.FlakeRegistry)
}

var CellsFrom = lazy.Of[string]{
	New: func() string {
		if s, err := getLocalCells(); err != nil {
			return "${cellsFrom}"
		} else {
			return s
		}
	},
}

// tprintf passed template string is formatted usign its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func tprintf(data interface{}, tmpl string) string {
	t := template.Must(template.New("tmp").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func getNix() (string, error) {
	nix, err := exec.LookPath("nix")
	if err != nil {
		return "", errors.New("You need to install 'nix' in order to use this tool")
	}
	return nix, nil
}

func getCurrentSystem() (string, error) {
	// detect the current system
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	currentSystem, err := exec.Command(
		nix, "eval", "--raw", "--impure", "--expr", "builtins.currentSystem",
	).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", err
	}
	currentSystemStr := string(currentSystem)
	return currentSystemStr, nil
}

func (r *FlakeRegistry) getCells() (string, error) {
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	cellsFrom, err := exec.Command(
		nix, "eval", "--raw", r.RefCellsFrom(),
	).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", err
	}
	return string(cellsFrom[:]), nil
}

func getLocalCells() (string, error) {
	// NB: has to create temporary var here: "Cannot take pointer off `LocalFlakeRegistry`"
	// sounds horrifyingly similar to C++'s rvalue/xvalue :)
	reg := LocalFlakeRegistry()
	return reg.getCells()
}
