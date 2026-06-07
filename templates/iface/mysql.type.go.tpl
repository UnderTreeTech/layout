{{- $short := (shortname .Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- end -}}

type {{ .Name }} interface {
    /********************Generate By XO: write method segment*********************/
    // Add{{ .Name }} add the {{ .Name }} to the database
    Add{{ .Name }}(ctx context.Context, {{$short}} *model.{{.Name}}) (err error)

    // Batch{{ .Name }} batch add {{ .Name }} to the database
    BatchAdd{{ .Name }}s(ctx context.Context, list []*model.{{.Name}}) (err error)

    // Edit{{ .Name }} edit {{ .Name }} in the database
    Edit{{ .Name }}(ctx context.Context, set map[string]interface{}, condition map[string]interface{}) (err error)

    // Delete{{ .Name }} delete {{ .Name }} from the database
    Delete{{ .Name }}(ctx context.Context, condition map[string]interface{}) (err error)

    /********************Generate By XO: read method segment**********************/
    // Find{{ .Name }} query one {{ .Name }} by condition from the database
    Find{{ .Name }}(ctx context.Context, condition map[string]interface{}) (res *model.{{ .Name }}, err error)

    // Find{{ .Name }}s all {{ .Name }} by condition from the database
    Find{{ .Name }}s(ctx context.Context, condition map[string]interface{}) (res []*model.{{ .Name }}, err error)

    // Count{{ .Name }} return count by condition from the database
    Count{{ .Name }}(ctx context.Context, condition map[string]interface{}) (num int, err error)

    // BatchCount{{ .Name }} groups records by groupKey and returns the count of each group.
    // The result is a slice of maps, where each map's key is the groupKey field value
    // and the value is the number of records in that group.
    // Example: groupKey="status" returns [{1: 10}, {2: 5}] meaning 10 records with status=1, 5 with status=2.
    BatchCount{{ .Name }}(ctx context.Context, groupKey string, condition map[string]interface{}) (res []map[string]int, err error)
}