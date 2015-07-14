# rsyncdiff
```bash
Difference verification tool of rsync command
Usage of rsyncdiff:
BuildDate:
rsyncdiff [OPTIONS] RSYNC-FROM RSYNC-TO

  -c, --context-diff=false    Produce a context format diff
  -e, --exclude=              rsync exclude option
  --exclude-from=             rsync exlude from option
  -l, --less=false            using less for output
  -r, --colordiff=false       using colordiff for output
  -t, --target-file=          Difference acquisition object file
  --unified-diff=true         Produce a unified format diff (default)
  -v, --vimdiff=false         Produce a vimdiff. Specify also t option
```

# TODO
- [ ] remote対応
