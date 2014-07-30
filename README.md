packer-provisioner-host-command
===============================

A packer provisioner that works much like the typical packer 'shell' provisioner,
except that it works on the *host* machine (i.e. the packer build machine) instead
of the guest/remote environment. The rationale for this is to have an automated step
to run something like serverspec tests against a packer built machine.

Note: most of this provisioner takes common work from the existing shell provisioner
in packer: https://github.com/mitchellh/packer/tree/master/provisioner/shell
