{{- $short := (shortname .Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
        {{ .Name }} {{ retype .Type }} `json:"{{ .Col.ColumnName }}"`
{{- end }}
}

// {{ .Name }}Columns defines and stores column names for '{{ $table }}'.
type {{ .Name }}Columns struct {
{{- range .Fields }}
    {{ .Name }} string
{{- end }}
}

// Get{{ .Name }}TableName returns the database table name for {{ .Name }}.
func Get{{  .Name  }}TableName() string {
    return "{{  .Table.TableName  }}"
}

// Get{{ .Name }}FieldString returns all column names as a comma-separated string (for SELECT).
func Get{{ .Name }}FieldString () string {
    return `{{ colnames .Fields }}`
}

// Get{{ .Name }}FieldStringSlice returns insert column names as a string slice (excludes auto-increment primary key).
func Get{{ .Name }}FieldStringSlice () []string {
    return strings.Split(`{{colnames .Fields .PrimaryKey.Name }}`,",")
}

// Get{{ .Name }}AddField returns the field values of {{ $short }} as a slice for INSERT operations.
func Get{{ .Name }}AddField ( {{ $short }}  *{{ .Name }}) []interface{} {
    return [] interface{} {
            		{{- range .Fields }}
            		    {{- if not .Col.IsPrimaryKey }}
                            {{ $short }}.{{ .Name }},
                        {{- end }}
                    {{- end }}
            }
}

// Get{{ .Name }}ScanField returns pointers to all fields of {{ $short }} for use with rows.Scan().
func Get{{ .Name }}ScanField ( {{ $short }}  *{{ .Name }}) []interface{} {
    return [] interface{} {
    		{{- range .Fields }}
                    &{{ $short }}.{{ .Name }},
            {{- end }}
    	}
}

// Get{{ .Name }}Columns returns a {{ .Name }}Columns struct with all column name mappings.
func Get{{ .Name }}Columns () {{ .Name }}Columns {
    return {{ .Name }}Columns {
    		{{- range .Fields }}
    		    {{ .Name }}: "{{ .Col.ColumnName }}",
    		{{- end }}
    	}
}
