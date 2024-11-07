package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	args := []string{"list", "-mod=readonly", "-m", "-u", "-json", "all"}
	cmd := exec.Command("go", args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("out:", outb.String())

	d := json.NewDecoder(strings.NewReader(outb.String()))
	for d.More() {
		var out struct {
			Path       string       `json:"Path,omitempty"`       // module path
			Version    string       `json:"Version,omitempty"`    // module version
			Versions   []string     `json:"Versions,omitempty"`   // available module versions
			Replace    *Module      `json:"Replace,omitempty"`    // replaced by this module
			Time       *time.Time   `json:"Time,omitempty"`       // time version was created
			Update     *Module      `json:"Update,omitempty"`     // available update (with -u)
			Main       bool         `json:"Main,omitempty"`       // is this the main module?
			Indirect   bool         `json:"Indirect,omitempty"`   // module is only indirectly needed by main module
			Dir        string       `json:"Dir,omitempty"`        // directory holding local copy of files, if any
			GoMod      string       `json:"GoMod,omitempty"`      // path to go.mod file describing module, if any
			GoVersion  string       `json:"GoVersion,omitempty"`  // go version used in module
			Retracted  []string     `json:"Retracted,omitempty"`  // retraction information, if any (with -retracted or -u)
			Deprecated string       `json:"Deprecated,omitempty"` // deprecation message, if any (with -u)
			Error      *ModuleError `json:"Error,omitempty"`      // error loading module
		}
		err := d.Decode(&out)
		if err != nil {
			panic(err)
		}
		if out.Update != nil {
			fmt.Println(out.Path, out.Version, "->", out.Update.Version)
		}
	}
}

type Module struct {
	Path       string       `json:"Path,omitempty"`       // module path
	Version    string       `json:"Version,omitempty"`    // module version
	Versions   []string     `json:"Versions,omitempty"`   // available module versions
	Replace    *Module      `json:"Replace,omitempty"`    // replaced by this module
	Time       *time.Time   `json:"Time,omitempty"`       // time version was created
	Update     *Module      `json:"Update,omitempty"`     // available update (with -u)
	Main       bool         `json:"Main,omitempty"`       // is this the main module?
	Indirect   bool         `json:"Indirect,omitempty"`   // module is only indirectly needed by main module
	Dir        string       `json:"Dir,omitempty"`        // directory holding local copy of files, if any
	GoMod      string       `json:"GoMod,omitempty"`      // path to go.mod file describing module, if any
	GoVersion  string       `json:"GoVersion,omitempty"`  // go version used in module
	Retracted  []string     `json:"Retracted,omitempty"`  // retraction information, if any (with -retracted or -u)
	Deprecated string       `json:"Deprecated,omitempty"` // deprecation message, if any (with -u)
	Error      *ModuleError `json:"Error,omitempty"`      // error loading module
}

type ModuleError struct {
	Err string `json:"Err,omitempty"` // the error itself
}
