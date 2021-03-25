# winclean

Looks for files older than a certain number hours (defaults to 72) in the Windows appdata/local/temp folder and removes them. Also looks within subfolders, checking file times within folders and removing them, and then removes the parent folder if it is empty.

Files that cannot be inspected or removed will show in the output.

Use the `--help` parameter for options.

`winclean.exe --help`
