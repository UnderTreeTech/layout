{{ $typeName := .Type.Name -}}
{{- $table := (schema .Schema .Type.Table.TableName) }}

// {{ .FuncName }}Iface defines the interface for {{ $typeName }} queries via index '{{ .Index.IndexName }}'.
// NOTE: embed this into the {{ $typeName }} DAO interface if needed, or add the method directly.
type {{ .FuncName }}Iface interface {
    // Find{{ .FuncName }} retrieves {{ if .Index.IsUnique }}a row{{ else }}rows{{ end }} from '{{ $table }}'
    // via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
    Find{{ .FuncName }}(ctx context.Context{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}res []*model.{{ $typeName }}, err error{{ else }}res *model.{{ $typeName }}, err error{{ end }})

    // Edit{{ .FuncName }} updates rows in '{{ $table }}'
    // via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
    Edit{{ .FuncName }}(ctx context.Context, set map[string]interface{}{{ goparamlist .Fields true true }}) (err error)

    // Delete{{ .FuncName }} deletes rows from '{{ $table }}'
    // via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
    Delete{{ .FuncName }}(ctx context.Context{{ goparamlist .Fields true true }}) (err error)
}
