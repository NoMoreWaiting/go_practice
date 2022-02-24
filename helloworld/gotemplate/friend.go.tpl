hello {{.Name}}
{{range .Emails}}
an email {{.}}
{{end}}
{{with .Friends}}
{{range .}}
an friend name is {{.Name}}
{{end}}
{{end}} 
