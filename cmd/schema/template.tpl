package {{.Package}} 

import "github.com/TRICERA-energy/sunspec"

var Schema = []sunspec.ModelDef{ {{range .Schema}}{{template "model" .}}{{end}} }

{{define "meta"}}{{if .Label}}
Label: "{{.Label}}",{{end}} {{if .Description}}
Description: "{{.Description}}",{{end}} {{if .Detail}}
Detail: "{{.Detail}}",{{end}} {{if .Notes}}
Notes: "{{.Notes}}",{{end}} {{if .Comments}} [ {{range .Comments}}"{{.}}",{{end}} ] ,{{end}}{{end}}

{{define "model"}} { {{if .Id}}
Id: {{.Id}},{{end}}{{template "meta" .}}
Group: sunspec.GroupDef{{template "group" .Group}}, } ,{{end}}

{{define "group"}} { {{if .Name}}
Name: "{{.Name}}",{{end}} {{if .Count}}
Count: {{printf "%#v" .Count}}{{end}} {{if .Atomic}}
Atomic: true,{{end}}{{template "meta" .}} {{if .Points}}
Points: []sunspec.PointDef{ {{range .Points}}{{template "point" .}},{{end}}
},{{end}} {{if .Groups}}
Groups: []sunspec.GroupDef{ {{range .Groups}}{{template "group" .}},{{end}}
},{{end}} } {{end}}

{{define "point"}} { {{if .Name}}
Name: "{{.Name}}",{{end}} {{if .Type}}
Type: "{{.Type}}",{{end}} {{if .Value}}
Value: {{printf "%#v" .Value}},{{end}} {{if .Count}}
Count: {{printf "%#v" .Count}},{{end}} {{if .Size}}
Size: {{.Size}},{{end}} {{if .ScaleFactor}}
ScaleFactor: {{printf "%#v" .ScaleFactor}},{{end}} {{if .Units}}
Units: {{.Units}},{{end}} {{if .Writable}}
Writable: true,{{end}} {{if .Mandatory}}
Mandatory: true,{{end}} {{if .Static}}
Static: true,{{end}}
} {{end}}
