package main

import (
	"testing"

	"github.com/mitchellh/packer/packer"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"commands": []interface{}{"foo", "bar"},
	}
}

func TestImplementsInterface(t *testing.T) {
	var p interface{}
	p = new(HostCommandProvisioner)

	if _, ok := p.(packer.Provisioner); !ok {
		t.Fatalf("HostCommandProvisioner is not a packer.Provisioner implementation")
	}
}

func TestPrepareAssignsDefaults(t *testing.T) {
	p := new(HostCommandProvisioner)
	config := map[string]interface{}{"command": "foo"}

	if err := p.Prepare(config); err != nil {
		t.Fatalf("err: %s", err)
	}

	if p.config.ExecuteCommand != "{{ .Command }}" {
		t.Errorf("Unexpected value for ExecuteCommand: %s", p.config.ExecuteCommand)
	}
}

func TestPrepareHandlesInvalidKeys(t *testing.T) {
	p := new(HostCommandProvisioner)
	config := map[string]interface{}{"command": "foo", "invalid_key": "bar"}

	if err := p.Prepare(config); err == nil {
		t.Fatal("Prepare with invalid key should have errored")
	}
}

func TestPrepareRejectsBothCommandAndCommands(t *testing.T) {
	p := new(HostCommandProvisioner)
	config := map[string]interface{}{"command": "foo", "commands": []string{"bar"}}

	if err := p.Prepare(config); err == nil {
		t.Fatal("Prepare with both command and commands should have errored")
	}
}

func TestPrepareWithCommandsSucceeds(t *testing.T) {
	p := new(HostCommandProvisioner)

	// Handle the case of no commands
	config := map[string]interface{}{"commands": []string{}}
	if err := p.Prepare(config); err == nil {
		t.Fatal("Prepare with empty commands should have errored")
	}

	// Handle the case of slice contents
	config = map[string]interface{}{"commands": []string{"bar"}}
	if err := p.Prepare(config); err != nil {
		t.Fatal("Prepare with commands should not have errored")
	}
}

func TestPrepareEnvironmentVars(t *testing.T) {
	config := map[string]interface{}{"commands": []string{"foo"}}

	// We should reject invalid enfironment variables
	p := new(HostCommandProvisioner)
	config["environment_vars"] = []string{"foo", "foo=bar", "=bar"}
	if err := p.Prepare(config); err == nil {
		t.Fatal("Prepare with invalid env vars should have errored")
	}

	// Ensure valid env vars
	p = new(HostCommandProvisioner)
	config["environment_vars"] = []string{"foo=bar", "baz="}
	if err := p.Prepare(config); err != nil {
		t.Fatal("Prepare with valid env vars should not have errored")
	}
}
