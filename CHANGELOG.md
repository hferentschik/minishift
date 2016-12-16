# Minishift Release Notes

# Version 1.0.0-beta.1 - 2016-10-19

* Issue [#193](https://github.com/minishift/minishift/issues/193) - bug - B2D ISO build uses latest boot2docker/boot2docker image for building
* Issue [#189](https://github.com/minishift/minishift/issues/189) - task - Standardize "Minishift"/"minishift" usage in the docs
* Issue [#186](https://github.com/minishift/minishift/issues/186) - bug - Building the default b2d ISO fails
* Issue [#171](https://github.com/minishift/minishift/issues/171) - feature - Release artifacts should be tar/zip bundles
* Issue [#154](https://github.com/minishift/minishift/issues/154) - task - Unused variable in Makefile
* Issue [#151](https://github.com/minishift/minishift/issues/151) - task - Remove the embedding of openshift binary into minishift
* Issue [#147](https://github.com/minishift/minishift/issues/147) - task - Remove the vendor directory and add it to .gitignore
* Issue [#145](https://github.com/minishift/minishift/issues/145) - task - Add AppVeyor CI
* Issue [#141](https://github.com/minishift/minishift/issues/141) - feature - Use 'oc cluster up' to start the OpenShift cluster
* Issue [#138](https://github.com/minishift/minishift/issues/138) - task - Review and update documentation
* Issue [#136](https://github.com/minishift/minishift/issues/136) - feature - Provide a custom provisioner for CentOS / RHEL based ISO images

Full [changelog](https://github.com/minishift/minishift/compare/v0.9.0...v1.0.0-beta.1)

# Version 0.9.0 - 2016-10-19

* [FEATURE] Upgrade to OpenShift v1.3.1 by default
* [FEATURE] Validate checksums of downloaded files
* [BUGFIX] Only download non-embedded version of OpenShift if required
* [FEATURE] Add minishift config view subcommand to be able to view current config

Full [changelog](https://github.com/minishift/minishift/compare/v0.8.0...v0.9.0)

# Version 0.8.0 - 2016-10-13

* [FEATURE] Default to 2 cpus & 2GB RAM
* [FEATURE] New flag to start subcommand to specify version of openshift to run
* [FEATURE] Config subcommand to persist all config options (CPU, memory, openshift version, etc)
* [FEATURE] Improved caching of downloaded files
* [FEATURE] Download progress bars

Full [changelog](https://github.com/minishift/minishift/compare/v0.7.1...v0.8.0)

# Version 0.7.1 - 2016-09-16

* [BUG] Fix OS default VM drivers

Full [changelog](https://github.com/minishift/minishift/compare/v0.7.0...v0.7.1)

# Version 0.7.0 - 2016-09-16

* [UPGRADE] OpenShift 1.3.0
* [FEATURE] Enable OpenShift registry by default
* [FEATURE] Support Kubernetes service annotation proposal
* [FEATURE] Show current version in update prompt
* [BUG] Fix wrong update file check

Full [changelog](https://github.com/minishift/minishift/compare/v0.6.0...v0.7.0)

# Version 0.6.0 - 2016-09-08
* Upgrade to OpenShift 1.3.0-rc1

# Version 0.5.0 - 2016-09-07
* [FEATURE] Enable host path provisioner
* [BREAKING] Rename VM to `minishift`
* [BUG] Fix xhyve hostname to `minishift`, rather than `boot2docker`
* [BUG] Ensure node IP is routeable
* [FEATURE] Reuse generated CA certificate
* [FEATURE] Ensure xhyve driver uses same IP on restarts
* [FEATURE] Add defaut insecure registry flag to include minishift service IP range
* [FEATURE] Allow environment variables to specify flags

# Version 0.3.2 - 2016-07-21
 * [BUG] Fix autoupdate checksums

# VERSION 0.3.1 - 2016-07-21
 * [BUG] Fix start command when running under xhyve on OS X

# Version 0.3.0 - 2016/07/18
 * BREAKING: Rename dashboard command to console
 * Add flag to pass extra Docker env vars to start command
 * Set router subdomain to <ip>.xip.io by default
 * EXPERIMENTAL: Auto-update of binaries
 * Build enhancements

# Version 0.2.1 - 2016/07/15
 * Enable all CORS origins for API server
 * Strip binary for smaller download
 * Build enhancements to check for valid cross compilation

## Version 0.2.0 - 2016/07/14
 * Changed API server port to 8443 to allow OpenShift router/ingress controllers to bind to 443 if required
 * Added a --disk-size flag to minishift start.
 * Fixed a bug regarding auth tokens not being reconfigured properly after VM restart
 * Added a new get-openshift-versions command, to get the available OpenShift versions so that users know what versions are available when trying to select the OpenShift version to use
 * Makefile Updates
 * Documentation Updates

## Version 0.1.1 - 2016/07/08
 * [BUG] Fix PATH problems preventing proper start up of OpenShift<Paste>

## Version 0.1.0 - 2016/07/07
 * Initial minishift  release.
