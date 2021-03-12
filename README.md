# winclean

Looks for files older than 48 hours in the Windows appdata/local/temp folder and removes them. Also looks within subfolders, checking file times within folders and removing them, and then removes the parent folder if it is empty.

Files that cannot be inspected or removed will show in the output.

For verbose output:
```
winclean.exe --verbose=true
```

Add `--dry=true` to only print actions without performing them.

