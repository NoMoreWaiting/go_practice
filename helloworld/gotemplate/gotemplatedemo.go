package gotemplate

import (
	"io"
	"text/template"
)

type Friend struct {
	Name string
}

type Person struct {
	Name    string
	Emails  []string
	Friends []*Friend
}

func TranslateTemplate(wr io.Writer, configPrefix string) {
	f1 := Friend{"f1"}
	f2 := Friend{"f2"}
	person := Person{"songyunxuan", []string{"123@163.com", "456@gmail.com"}, []*Friend{&f1, &f2}}

	tpl := template.Must(template.ParseFiles(configPrefix + "friend.go.tpl"))

	tpl.Execute(wr, person)
}
