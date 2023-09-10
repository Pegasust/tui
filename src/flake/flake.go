package flake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/hymkor/go-lazy"
)

var (
	BrandedRegistry = "__std" // keep for now for historic reasons
)

type PaisanoRegistry struct {
	FlakeRef        string `regroup:"flake_ref,optional"`
	Registry        string `regroup:"registry,optional"`
	PaisanoRegistry string
}

func LocalPaisanoRegistry() PaisanoRegistry {
	return PaisanoRegistry{
		FlakeRef:        "./",
		Registry:        BrandedRegistry,
		PaisanoRegistry: fmt.Sprintf("%s#%s", "./", BrandedRegistry),
	}
}

func (r *PaisanoRegistry) PrefetchFlakeRef() (string, error) {
	nix, err := getNix()
	if err != nil {
		return "", err
	}
	prefetchOut, err := exec.Command(
		nix, "flake", "prefetch", "--json", r.FlakeRef,
	).Output()
	if err != nil {
		fmterr := err
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmterr = fmt.Errorf("%w, stderr:\n%s", exitErr, exitErr.Stderr)
		}
		return "", fmt.Errorf("failed to prefetch, is this a flake? %v", fmterr)
	}

	hashStore := struct {
		Hash      string `json:"hash"`
		StorePath string `json:"storePath"`
	}{
		Hash:      "",
		StorePath: "",
	}
	err = json.Unmarshal(prefetchOut, &hashStore)
	if err != nil {
		return "", fmt.Errorf("failed to parse prefetch path, programmer's err: %v", err)
	}
	return hashStore.StorePath, nil
}

func (r *PaisanoRegistry) PrefetchPaisanoReg() (string, error) {
	localizedFlakeRef, err := r.PrefetchFlakeRef()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s#%s", localizedFlakeRef, r.Registry), nil
}

func (r *PaisanoRegistry) PaisanoReg() string {
	flakeReg, err := r.PrefetchPaisanoReg()
	if err != nil {
		// NOTE: warn?
		flakeReg = r.PaisanoRegistry
	}
	return flakeReg
}

func (r *PaisanoRegistry) InitPaisanoRef(system string) string {
	return fmt.Sprintf("%s.init.%s", r.PaisanoReg(), system)
}

func (r *PaisanoRegistry) RefCellsFrom() string {
	return fmt.Sprintf("%s.cellsFrom", r.PaisanoReg())
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

func (r *PaisanoRegistry) getCells() (string, error) {
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
	reg := LocalPaisanoRegistry()
	return reg.getCells()
}
