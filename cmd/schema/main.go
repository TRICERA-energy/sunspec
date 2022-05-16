package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/TRICERA-energy/sunspec"
)

func main() {
	pkg := flag.String("p", "schema", "")
	out := flag.String("o", "./schema", "")
	dat := flag.String("d", "", "")

	flag.Parse()

	target := fmt.Sprintf("%v.go", *out)

	var schema []sunspec.ModelDef
	must(json.Unmarshal([]byte(*dat), &schema))

	tpl := template.Must(template.ParseFiles("./template.tpl"))

	defer exec.Command("go", "fmt", target).Run()

	f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	must(err)
	defer f.Close()

	must(tpl.Execute(f, struct {
		Package string
		Schema  []sunspec.ModelDef
	}{*pkg, schema}))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
