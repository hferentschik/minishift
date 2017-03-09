# Troubleshooting

This section contains solutions to common problems that you might encounter
while using Minishift.

<!-- MarkdownTOC -->

- [KVM driver](#kvm-driver)
  - [Error creating new host: dial tcp: missing address](#error-creating-new-host-dial-tcp-missing-address)
  - [Failed to connect socket to '/var/run/libvirt/virtlogd-sock'](#failed-to-connect-socket-to-varrunlibvirtvirtlogd-sock)
  - [Error starting the VM: ... operation failed: domain 'minishift' already exists ...](#error-starting-the-vm--operation-failed-domain-minishift-already-exists-)
- [xhyve driver](#xhyve-driver)
  - [Error: could not create vmnet interface, permission denied or no entitlement](#error-could-not-create-vmnet-interface-permission-denied-or-no-entitlement)
- [VirtualBox driver](#virtualbox-driver)
  - [Error: getting state for host: machine does not exist](#error-getting-state-for-host-machine-does-not-exist)
- [Users and authentication](#users-authentication)
  - [Some special characters cause passwords to fail](#some-special-characters-cause-passwords-to-fail)
- [Snapshots with virsh](#snapshots-with-virsh)
  - [`minishift delete` fails to undefine snapshots of running instances](#minishift-delete-fails-to-undefine-snapshots-of-running-instances)

<!-- /MarkdownTOC -->


<a name="kvm-driver"></a>
## KVM driver

<a name="error-creating-new-host-dial-tcp-missing-address"></a>
### Error creating new host: dial tcp: missing address

The problem is likely to be that the `libvirtd` service is not running, you can check it with

```
systemctl status libvirtd
```

If `libvirtd` is not running, start and enable it to start on boot:

```
systemctl start libvirtd
systemctl enable libvirtd
```

<a name="failed-to-connect-socket-to-varrunlibvirtvirtlogd-sock"></a>
### Failed to connect socket to '/var/run/libvirt/virtlogd-sock'

The problem is likely to be that the `virtlogd` service is not running, you can check it with

```
systemctl status virtlogd
```

If `virtlogd` is not running, start and enable it to start on boot:

```
systemctl start virtlogd
systemctl enable virtlogd
```

<a name="error-starting-the-vm--operation-failed-domain-minishift-already-exists-"></a>
### Error starting the VM: ... operation failed: domain 'minishift' already exists ...

Check for existing VMs and remove them:

```
sudo virsh list --all
sudo virsh destroy minishift
sudo virsh undefine minishift
```

<a name="xhyve-driver"></a>
## xhyve driver

<a name="error-could-not-create-vmnet-interface-permission-denied-or-no-entitlement"></a>
### Error: could not create vmnet interface, permission denied or no entitlement

The problem is likely to be that the xhyve driver is not able to clean up
vmnet when a VM is removed. vmnet.framework decides the IP based on following files:

* _/var/db/dhcpd_leases_
* _/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist_

Reset the `minishift` specific IP database and make sure you remove `minishift`
entry section from `dhcpd_leases` file. Finally, reboot your system.

    {
      ip_address=192.168.64.2
      hw_address=1,2:51:8:22:87:a6
      identifier=1,2:51:8:22:87:a6
      lease=0x585e6e70
      name=minishift
    }

**Note:** You can completely reset IP database by removing both the files
manually but this is very **risky**.

<a name="virtualbox-driver"></a>
## VirtualBox driver

<a name="error-getting-state-for-host-machine-does-not-exist"></a>
### Error: getting state for host: machine does not exist

If you use Windows, ensure that you used the `--vm-driver virtualbox` flag with the `minishift start` command. Alternatively, the problem is likely to be an outdated version
of Virtual Box.

It is recommended to use `Virtualbox >= 5.1.12` to avoid this issue.

<a name="users-authentication"></a>
## Users and authentication

<a name="some-special-characters-cause-passwords-to-fail"></a>
### Some special characters cause passwords to fail

Depending on your operating system and shell environment, certain special characters
can trigger variable interpolation and therefore cause passwords to fail.

Workaround: When creating and entering passwords, wrap the string with single quotes in
the following format: '&lt;password>'

<a name="snapshots-with-virsh"></a>
## Snapshots with virsh

<a name="minishift-delete-fails-to-undefine-snapshots-of-running-instances"></a>

### `minishift delete` fails to undefine snapshots of running instances, made using virsh on KVM/libvirt.

If you use virsh on KVM/libvirt to create snapshots in your development workflow, using `minishift delete` to delete the snapshots, along with the VM, returns an error:

    $ minishift delete
    Deleting the Minishift VM...
    Error deleting the VM:  [Code-55] [Domain-10] Requested operation is not valid: cannot delete inactive domain with 4 snapshots

Workaround: The snapshots are stored in `~/.minishift/machines`, but the definitions are stored in `var/lib/libvirt/qemu/snapshot/minishift`.

To delete the snapshots you need to:

1. Delete the definitions using:

        $ sudo virsh snapshot-delete --metadata minishift <snapshot-name>

1. Undefine the Minishift domain using:

        $ sudo virsh undefine minishift

**Note:** In case the above step does not resolve the issue, you can also use the following command to delete the snapshots:

    $ rm -rf ~/.minishift/machines

You can now do `minishft delete` to delete the VM and restart Minishift.

It is recommended to avoid use of metadata when you create snapshots as follows:

    $ sudo virsh snapshot-create-as --domain vm1 overlay1 --diskspec vda,file=/export/overlay1.qcow2 --disk-only --atomic --no-metadata
