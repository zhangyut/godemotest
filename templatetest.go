package main

import (
	"fmt"
	"os"
	"text/template"
)

const view = `
view {{.Name}} {
    match-clients { {{range .Ips}} {{.}}; {{end}} };
{{range .Zones}}
	zone {{.Name}} {
		type forward;
		forward {{.Type}};
	};
{{end}}
};
`

type Zone struct {
	Name string
	Type string
}

type View struct {
	Name  string
	Ips   []string
	Zones []Zone
}

func main() {
	zones := []Zone{
		{
			Name: "zone1",
			Type: "type1",
		},
		{
			Name: "zone2",
			Type: "type2",
		},
		{
			Name: "zone3",
			Type: "type3",
		},
	}
	views := []View{
		{
			Name:  "view1",
			Zones: zones,
			Ips:   []string{"1.2.3.4", "5.6.7.8/24"},
		},
		{
			Name:  "view2",
			Zones: zones,
		},
	}
	t := template.Must(template.New("views").Parse(view))

	for _, r := range views {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			fmt.Println("executing template:", err)
		}
	}
}
