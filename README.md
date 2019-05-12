# troubleshootbrokentimemachinevolume
troubleshoot broken timemachine volume

## Draft 1

Collect path and inode information into sqlite3 database

```shell
sudo ./troubleshootbrokentimemachinevolume foo.db /Volume/FooTimeMachine
```

## Caution

macOS Mojave by default prevents programs run from within Terminal.app from doing exciting things like listing directory contents of Time Machine backups.
Terminal.app needs "Full Disk Access" in System Preferences, Security & Privacy.
See [Fix Terminal “Operation not permitted” Error in MacOS Mojave](http://osxdaily.com/2018/10/09/fix-operation-not-permitted-terminal-error-macos/) for instructions.
