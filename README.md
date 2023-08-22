## MVP is mv with patterns

Usacases:

```bash
$ ls
log_1.json  log_2.json warn_3.json error_3.json README.md

$ mvp \$1_\$2.json \$2_\$1.json
$ ls
1_log.json  2_log.json 3_warn.json 3_error.json README.md`

```

Use `$<number>` like `$1` `$2` etc to match patterns and use these matched patterns in the new file names.
If you have files like `log_2023-08-22_dev.json warn_2023-08-22_prod.json warn_2023-08-23_stage.json`. They clearly
have this pattern `<log level>_<date>_<environment>.json`. Use `mvp` to match these.
The pattern `\$1_\$2_\$3.json` will replace `$1` with the **log level**, `$2` with the **date** and `$3` with the
**environment**. We might be interested in a new naming. Something like `<environment>_<log level>_<date>.json`

Before:

```bash
## Before mvp
$ ls
log_2023-08-22_dev.json warn_2023-08-22_prod.json warn_2023-08-23_stage.json

$ mvp \$1_\$2_\$3.json \$3_\$1_\$2.json

## After mvp
$ ls
dev_log_2023-08-22.json prod_warn_2023-08-22.json stage_warn_2023-08-23.json
```
