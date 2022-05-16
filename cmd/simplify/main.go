package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/TRICERA-energy/sunspec"
)

func main() {
	in := flag.String("i", "", "")
	flag.Parse()

	var schema []sunspec.ModelDef
	must(json.Unmarshal([]byte(*in), &schema))

	for i := range schema {
		schema[i].Simplify()
	}

	d, err := json.MarshalIndent(schema, "", "\t")
	must(err)

	fmt.Fprintf(os.Stdout, "%v", string(d))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
