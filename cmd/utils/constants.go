package utils

import "path"

var PathToTemplates = path.Join("./", "cmd", "templates")

const MinIdentityLen = 8
const MaxIdentityLen = 32
const MinNameLen = 1
const MaxNameLen = 255
const EmailRgx = `^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`
