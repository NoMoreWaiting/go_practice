{{ /* goland 需要将 go template 项目置于 gopath 中， 并且需要额外的手工引导标记， 不方便且不通用， 暂时放弃 */ }}

hello {{.Name}}
{{range .Emails}}
an email {{.}}
{{end}}
{{with .Friends}}
{{range .}}
an friend name is {{.Name}}
{{end}}
{{end}} 
