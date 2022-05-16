package main

import (
	"encoding/json"
	"flag"
	"os"
	"text/template"

	"github.com/TRICERA-energy/sunspec"
)

func main() {
	pkg := flag.String("p", "schema", "")
	dat := flag.String("d", "", "")

	flag.Parse()

	var schema []sunspec.ModelDef
	must(json.Unmarshal([]byte(*dat), &schema))

	tpl := template.Must(template.ParseFiles("./template.tpl"))

	must(tpl.Execute(os.Stdout, struct {
		Package string
		Schema  []sunspec.ModelDef
	}{*pkg, schema}))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
