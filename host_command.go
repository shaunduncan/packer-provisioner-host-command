// This provisioner works much like the typical packer 'shell' provisioner,
// except that it works on the *host* machine (i.e. the packer build machine)
// instead of the guest. The rationale for this is to have an automated step
// to run something like serverspec tests against a packer built machine
//
// Note: most of what is contained here mirrors nearly verbatim the source of
// the shell provisioner in packer: https://github.com/mitchellh/packer/tree/master/provisioner/shell
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/flynn/go-shlex"
	"github.com/mitchellh/packer/common"
        "github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/packer/plugin"
        "github.com/mitchellh/packer/template/interpolate"
)

type ExecuteCommandTemplate struct {
	Command string
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// Path of single command to execute
	Command string

	// Array of multiple commands to execute
	Commands []string

	// Environment variables that will be added prior to executing
	Vars []string `mapstructure:"environment_vars"`

	// The string used to execute a command. The command being executed
	// is denoted as {{ .Command }}
	ExecuteCommand string `mapstructure:"execute_command"`

	// The configuration template
        ctx interpolate.Context
}

// A wrapper for stdout/stderr to use ui.Say and ui.Error respectively
type CommandWriter struct {
	WriteFunc func(string)
}

func (w CommandWriter) Write(p []byte) (n int, err error) {
	w.WriteFunc(string(p))
	return len(p), nil
}

type HostCommandProvisioner struct {
	config Config
}

func (p *HostCommandProvisioner) Prepare(raw ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
                Interpolate: true,
                InterpolateFilter: &interpolate.RenderFilter{
                        Exclude: []string{
                                "execute_command",
                        },
               },
        }, raw...)

	if err != nil {
		return err
	}

        var errs *packer.MultiError

	if p.config.ExecuteCommand == "" {
		p.config.ExecuteCommand = "{{ .Command }}"
	}

	if p.config.Commands == nil {
		p.config.Commands = make([]string, 0)
	}

	if p.config.Vars == nil {
		p.config.Vars = make([]string, 0)
	}

	if p.config.Command != "" && len(p.config.Commands) > 0 {
		errs = packer.MultiErrorAppend(errs,
			errors.New("You cannot specify both command and commands"))
	}

	if p.config.Command != "" {
		p.config.Commands = []string{p.config.Command}
	}

	// Sanity check
	if len(p.config.Commands) == 0 {
		errs = packer.MultiErrorAppend(errs, errors.New("You must specify a command"))
	}

	// Check for bad env vars (i.e. '=foo' and 'foobar')
	for _, keyValStr := range p.config.Vars {
		keyValPair := strings.SplitN(keyValStr, "=", 2)
		if len(keyValPair) != 2 || keyValPair[0] == "" {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("Environment variable not in 'key=value' format: %s", keyValStr))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *HostCommandProvisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
	// Build our variables up by adding in the build name and builder type
	env := make([]string, len(p.config.Vars)+2)
	env[0] = "PACKER_BUILD_NAME=" + p.config.PackerBuildName
	env[1] = "PACKER_BUILDER_TYPE=" + p.config.PackerBuilderType
	copy(env[2:], p.config.Vars)

	// Run commands directly on the host machine
	for _, cmdStr := range p.config.Commands {
		ui.Say(fmt.Sprintf("Running host command: %s", cmdStr))

		// Compile the full command string
                p.config.ctx.Data = &ExecuteCommandTemplate{
                	Command: cmdStr,
                }

                command, err := interpolate.Render(p.config.ExecuteCommand, &p.config.ctx)

		if err != nil {
			return fmt.Errorf("Error processing command '%s': %s", cmdStr, err)
		}

		// We have to split the command string into the actual command and its arguments
		parts, err := shlex.Split(command)

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Env = env

		// Redirect stdout/stderr to the ui writer
		cmd.Stdout = CommandWriter{WriteFunc: ui.Say}
		cmd.Stderr = CommandWriter{WriteFunc: ui.Error}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Error running command '%s': %s", cmdStr, err)
		}
	}

	return nil
}

func (p *HostCommandProvisioner) Cancel() {
	// We're not talking to the guest machine with this provisioner, so just quit
	os.Exit(0)
}

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterProvisioner(new(HostCommandProvisioner))
	server.Serve()
}
