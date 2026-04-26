package plymouth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Descriptor struct {
	Name        string
	Description string
	ModuleName  string
	ImageDir    string
	ScriptFile  string
}

func Parse(path string) (*Descriptor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := &Descriptor{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "Name":
			d.Name = val
		case "Description":
			d.Description = val
		case "ModuleName":
			d.ModuleName = val
		case "ImageDir":
			d.ImageDir = val
		case "ScriptFile":
			d.ScriptFile = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Descriptor) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("missing Name")
	}
	if d.ModuleName != "script" {
		return fmt.Errorf("ModuleName must be script, got %q", d.ModuleName)
	}
	if d.ImageDir == "" {
		return fmt.Errorf("missing ImageDir")
	}
	if d.ScriptFile == "" {
		return fmt.Errorf("missing ScriptFile")
	}
	return nil
}
