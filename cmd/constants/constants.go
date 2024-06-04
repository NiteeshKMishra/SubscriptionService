package constants

import (
	"io/fs"
	"path"
)

var PathToTemplates = path.Join("./", "cmd", "templates")
var PathToAssets = path.Join("./", "assets")
var DirPermission fs.FileMode = 0777

const MinIdentityLen = 8
const MaxIdentityLen = 32
const MinNameLen = 1
const MaxNameLen = 255
const EmailRgx = `^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
