// Get{{ .FuncName }}IndexFields returns the column names for index '{{ .Index.IndexName }}'.
func Get{{ .FuncName }}IndexFields() []string {
    return []string{ {{ range $i, $f := .Fields }}{{if $i}}, {{end}}"{{ $f.Col.ColumnName }}"{{ end }} }
}

