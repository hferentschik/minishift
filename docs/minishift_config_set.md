## minishift config set

Sets an individual value in a minishift config file

### Synopsis


Sets the PROPERTY_NAME config value to PROPERTY_VALUE
	These values can be overwritten by flags or environment variables at runtime.

```
minishift config set PROPERTY_NAME PROPERTY_VALUE
```

### Options inherited from parent commands

```
      --alsologtostderr value          log to standard error as well as files
      --disable-update-notification    Whether to disable VM update check.
      --log-flush-frequency duration   Maximum number of seconds between log flushes (default 5s)
      --log_backtrace_at value         when logging hits line file:N, emit a stack trace (default :0)
      --log_dir value                  If non-empty, write log files in this directory
      --logtostderr value              log to standard error instead of files
      --password string                Password to register Virtual Machine
      --show-libmachine-logs           Whether or not to show logs from libmachine.
      --stderrthreshold value          logs at or above this threshold go to stderr (default 2)
      --username string                Username to register Virtual Machine
  -v, --v value                        log level for V logs
      --vmodule value                  comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [minishift config](minishift_config.md)	 - Modify minishift config

