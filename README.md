## MVP is mv with patterns

Usacases:

```
$ ls
log_1_info.json  log_2_info.json log_2_info.json README.md

$ mvp log_$_info.json log_$_warn.json # $ is (1,2,3)
$ ls
info_1_warn.json  debug_2_warn.json error_3_warn.json README.md

```

```
$ ls
log_1_info.json  log_2_debug.json log_2_error.json README.md

$ mvp log_$1_$2.json $2_$1_log.json # $1 matches (1,2,3); $2 matches (info,debug,error)
$ ls
info_1_log.json  debug_2_log.json error_3_log.json README.md
```

Use `$` or numbered `$1` `$2` etc to match patterns and use these matched patterns in the new file names
