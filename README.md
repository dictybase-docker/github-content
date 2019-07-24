# github-content
cli to download modified files from github commit

### Command line
```
Usage:
  github-content [flags]

Flags:
  -c, --commit-payload string   commit data that is received from GitHub after triggered by a push event[required]
      --doc                     generate markdown documentation
  -p, --file-extension string   file extension that will be screened in the commit payload (default "obo")
  -f, --folder string           output folder[required]
  -h, --help                    help for github-content
      --log-file string         file for log output other than standard output, written to a temp folder by default
      --log-format string       format of the logging out, either of json or text (default "json")
      --log-level string        log level for the application (default "error")
  -o, --owner string            github repository owner[required]
  -r, --repository string       github repository name[required]
```
required flag(s) "commit-payload", "folder", "owner", "repository"
