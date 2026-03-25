package opencli_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/block/opencli-go"
)

func TestMarshalAlwaysWritesVersion(t *testing.T) {
	doc := opencli.Document{
		OpenCLI: "something-else",
		Info:    opencli.CliInfo{Version: "1.0.0"},
		Command: opencli.Command{Name: "myapp"},
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if got := raw["opencli"]; got != opencli.Version {
		t.Errorf("opencli = %q, want %q", got, opencli.Version)
	}
}

func TestMarshalEmptyOpenCLI(t *testing.T) {
	doc := opencli.Document{
		Info:    opencli.CliInfo{Version: "1.0.0"},
		Command: opencli.Command{Name: "myapp"},
	}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if got := raw["opencli"]; got != opencli.Version {
		t.Errorf("opencli = %q, want %q", got, opencli.Version)
	}
}

func TestUnmarshalMissingOpenCLI(t *testing.T) {
	input := `{"name":"myapp","info":{"version":"1.0.0"}}`
	var doc opencli.Document
	err := json.Unmarshal([]byte(input), &doc)
	if err == nil {
		t.Fatal("expected error for missing opencli field")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("error = %q, want it to mention 'missing'", err.Error())
	}
}

func TestUnmarshalUnsupportedVersion(t *testing.T) {
	input := `{"opencli":"9.9","name":"myapp","info":{"version":"1.0.0"}}`
	var doc opencli.Document
	err := json.Unmarshal([]byte(input), &doc)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
	if !strings.Contains(err.Error(), "9.9") {
		t.Errorf("error = %q, want it to mention '9.9'", err.Error())
	}
}

func TestUnmarshalValidDocument(t *testing.T) {
	input := `{
		"opencli": "0.1-block.1",
		"name": "myapp",
		"info": {"version": "2.0.0"},
		"commands": [
			{"name": "sub", "options": [{"name": "--verbose"}]}
		],
		"options": [{"name": "--help", "aliases": ["-h"]}]
	}`
	var doc opencli.Document
	if err := json.Unmarshal([]byte(input), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.OpenCLI != "0.1-block.1" {
		t.Errorf("OpenCLI = %q, want %q", doc.OpenCLI, "0.1-block.1")
	}
	if doc.Name != "myapp" {
		t.Errorf("Name = %q, want %q", doc.Name, "myapp")
	}
	if doc.Info.Version != "2.0.0" {
		t.Errorf("Info.Version = %q, want %q", doc.Info.Version, "2.0.0")
	}
	if len(doc.Commands) != 1 || doc.Commands[0].Name != "sub" {
		t.Errorf("Commands = %v, want one command named 'sub'", doc.Commands)
	}
	if len(doc.Options) != 1 || doc.Options[0].Name != "--help" {
		t.Errorf("Options = %v, want one option named '--help'", doc.Options)
	}
}

func TestRoundTrip(t *testing.T) {
	original := opencli.Document{
		Info: opencli.CliInfo{Version: "1.0.0"},
		Command: opencli.Command{
			Name:    "myapp",
			Aliases: []string{"ma"},
			Commands: []opencli.Command{
				{Name: "serve"},
			},
			Options: []opencli.Option{
				{Name: "--port", Aliases: []string{"-p"}},
			},
		},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded opencli.Document
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Name != "myapp" {
		t.Errorf("Name = %q, want %q", decoded.Name, "myapp")
	}
	if len(decoded.Aliases) != 1 || decoded.Aliases[0] != "ma" {
		t.Errorf("Aliases = %v, want [ma]", decoded.Aliases)
	}
	if len(decoded.Commands) != 1 || decoded.Commands[0].Name != "serve" {
		t.Errorf("Commands round-trip failed")
	}
	if len(decoded.Options) != 1 || decoded.Options[0].Name != "--port" {
		t.Errorf("Options round-trip failed")
	}
}

func TestArityDefaultsWhenOmitted(t *testing.T) {
	input := `{"name":"arg1"}`
	var arg opencli.Argument
	if err := json.Unmarshal([]byte(input), &arg); err != nil {
		t.Fatal(err)
	}
	if arg.Arity.Minimum != 1 || arg.Arity.Maximum != 1 {
		t.Errorf("Arity = {%d, %d}, want {1, 1}", arg.Arity.Minimum, arg.Arity.Maximum)
	}
}

func TestArityDefaultsPartialOverride(t *testing.T) {
	input := `{"name":"arg1","arity":{"minimum":0}}`
	var arg opencli.Argument
	if err := json.Unmarshal([]byte(input), &arg); err != nil {
		t.Fatal(err)
	}
	if arg.Arity.Minimum != 0 {
		t.Errorf("Arity.Minimum = %d, want 0", arg.Arity.Minimum)
	}
	if arg.Arity.Maximum != 1 {
		t.Errorf("Arity.Maximum = %d, want 1", arg.Arity.Maximum)
	}
}
