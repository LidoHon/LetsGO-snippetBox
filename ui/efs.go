package ui

/* note
If a path is to a directory, then all files in that directory are
recursively embedded, except for files with names that begin with
. or _. If you want to include these files you should use the all:
prefix, like go:embed "all:static".
*/


import (
	"embed"
)
//go:embed "html" "static"
var Files embed.FS