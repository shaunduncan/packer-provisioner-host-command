packer-provisioner-host-command
===============================

A packer provisioner that works much like the typical packer 'shell' provisioner,
except that it works on the *host* machine (i.e. the packer build machine) instead
of the guest/remote environment. The rationale for this is to have an automated step
to run something like serverspec tests against a packer built machine.

Note: most of this provisioner takes common work from the existing shell provisioner
in packer: https://github.com/mitchellh/packer/tree/master/provisioner/shell


Usage
-----

Usage of this provisioner is similar to the ``shell`` provisioner. Except, the type
is ``host-command`` and fields are either ``command`` or ``commands``. For example:

```json
{
  "type": "host-command",
  "command": "ls -la"
}
```

Or multiple commands, run sequentially:

```json
{
  "type": "host-command",
  "commands": [
    "whoami",
    "pwd",
    "/path/to/script.sh"
  ]
}
```

Both ``command`` and ``commands`` are treated as [configuration templates](http://www.packer.io/docs/templates/configuration-templates.html).
The full list of available options for this provisioner are:

* **command**: A string of the command to run
* **commands**: A list of strings of commands to run. These are run sequentially.
* **environment_vars**: A list of strings for environment variables to set before
  running a command. These values must be valid syntax (i.e. "foo=bar" or "baz=")
* **execute_command**: The template for how a command is run. This is a configuration
  template, and by default is "{{ .Command }}". The template variable "{{ .Command }}"
  must be present.


Building and Running Tests
--------------------------

You will need Go installed, then run ``make``, which will perform a ``clean``,
``deps``, and ``build``.

To run the tests: ``make test``


License
-------

The MIT License (MIT)

Copyright (c) 2014 Shaun Duncan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
