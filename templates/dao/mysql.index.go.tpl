{{- $short := (shortname .Type.Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Type.Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- end -}}

/********************Generate By XO: index query method segment*********************/
// Find{{ .FuncName }} retrieves {{ if .Index.IsUnique }}a row{{ else }}rows{{ end }} from '{{ $table }}'
// via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
func (d *dao) Find{{ .FuncName }}(ctx context.Context{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}res []*model.{{ .Type.Name }}, err error{{ else }}res *model.{{ .Type.Name }}, err error{{ end }}) {
    // init build with index condition
    build := squirrel.Select(model.Get{{ .Type.Name }}FieldString()).
        From(model.Get{{ .Type.Name }}TableName()).
        Where("{{ colnamesquery .Fields " AND " }}"{{ goparamlist .Fields true false }})

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
        return
    }

    // run query
    log.Debug(ctx, "Find{{ .FuncName }}", log.String("sql", fmt.Sprint(sqlStr, args)))
{{- if .Index.IsUnique }}
    res = &model.{{ .Type.Name }}{}
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        err = d.db.QueryRow(ctx, sqlStr, args...).Scan(model.Get{{ .Type.Name }}ScanField(res)...)
    } else {
        err = tx.QueryRow(sqlStr, args...).Scan(model.Get{{ .Type.Name }}ScanField(res)...)
    }
    return
{{- else }}
    var rows *sql.Rows
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        rows, err = d.db.Query(ctx, sqlStr, args...)
    } else {
        rows, err = tx.Query(sqlStr, args...)
    }

    if err != nil {
        return
    }
    defer rows.Close()

    // load results
    res = make([]*model.{{ .Type.Name }}, 0)
    for rows.Next() {
        // scan rows
        {{ $short }} := &model.{{ .Type.Name }}{}
        err = rows.Scan(model.Get{{ .Type.Name }}ScanField({{ $short }})...)
        if err != nil {
            return
        }

        res = append(res, {{ $short }})
    }
    return
{{- end }}
}

// Edit{{ .FuncName }} updates rows in '{{ $table }}'
// via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
func (d *dao) Edit{{ .FuncName }}(ctx context.Context, set map[string]interface{}{{ goparamlist .Fields true true }}) (err error) {
    // init build with index condition
    build := squirrel.Update(model.Get{{ .Type.Name }}TableName()).
        SetMap(set).
        Where("{{ colnamesquery .Fields " AND " }}"{{ goparamlist .Fields true false }})

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
        return
    }

    // run query
    log.Debug(ctx, "Edit{{ .FuncName }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Delete{{ .FuncName }} deletes rows from '{{ $table }}'
// via index '{{ .Index.IndexName }}'. Index fields: {{ colnames .Fields }}.
func (d *dao) Delete{{ .FuncName }}(ctx context.Context{{ goparamlist .Fields true true }}) (err error) {
    // init build with index condition
    build := squirrel.Delete(model.Get{{ .Type.Name }}TableName()).
        Where("{{ colnamesquery .Fields " AND " }}"{{ goparamlist .Fields true false }})

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
        return
    }

    // run query
    log.Debug(ctx, "Delete{{ .FuncName }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

