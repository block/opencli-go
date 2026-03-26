package opencli

import (
	"encoding/json"
	"errors"
	"fmt"

	current "github.com/block/opencli-go/internal/v0_1_block_1"
)

// Version is the OpenCLI spec version implemented by this module.
const Version = "0.1-block.1"

// Public type aliases for the current version.
type (
	Argument    = current.Argument
	Arity       = current.Arity
	CliInfo     = current.CliInfo
	Command     = current.Command
	Contact     = current.Contact
	Conventions = current.Conventions
	Email       = current.Email
	ExitCode    = current.ExitCode
	License     = current.License
	Metadata    = current.Metadata
	Option      = current.Option
)

// Document is the root of an OpenCLI description.
// It embeds Command and adds the required opencli version and info fields.
type Document struct {
	// OpenCLI is the spec version. Populated on unmarshal; always written as Version on marshal.
	OpenCLI     string       `json:"opencli"`
	Info        CliInfo      `json:"info"`
	Conventions *Conventions `json:"conventions,omitempty"`
	Command
}

// MarshalJSON always writes Version as the "opencli" field.
func (d Document) MarshalJSON() ([]byte, error) {
	type document Document
	tmp := document(d)
	tmp.OpenCLI = Version
	return json.Marshal(tmp)
}

// UnmarshalJSON parses the "opencli" field first, then dispatches
// to version-specific parsing. Returns an error if the field is
// missing or the version is unrecognized.
func (d *Document) UnmarshalJSON(data []byte) error {
	var envelope struct {
		OpenCLI string `json:"opencli"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if envelope.OpenCLI == "" {
		return errors.New("opencli: missing required \"opencli\" field")
	}
	switch envelope.OpenCLI {
	case "0.1-block.1":
		// Unmarshal Document-level fields (opencli, info, conventions)
		// separately from the embedded Command, because Command has its
		// own UnmarshalJSON which would consume the entire input.
		var docFields struct {
			OpenCLI     string       `json:"opencli"`
			Info        CliInfo      `json:"info"`
			Conventions *Conventions `json:"conventions,omitempty"`
		}
		if err := json.Unmarshal(data, &docFields); err != nil {
			return err
		}
		if err := json.Unmarshal(data, &d.Command); err != nil {
			return err
		}
		d.OpenCLI = docFields.OpenCLI
		d.Info = docFields.Info
		d.Conventions = docFields.Conventions
		return nil
	default:
		return fmt.Errorf("opencli: unsupported version %q", envelope.OpenCLI)
	}
}
