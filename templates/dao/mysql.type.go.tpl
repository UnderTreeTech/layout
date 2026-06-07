{{- $short := (shortname .Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- end -}}

/********************Generate By XO: write method segment*********************/
// Add{{ .Name }} add the {{ .Name }} to the database
func (d *dao) Add{{ .Name }}(ctx context.Context, {{$short}} *model.{{.Name}}) (err error) {
    // sql insert query, primary key provided by autoincrement
	sqlStr, args, err := squirrel.
        		Insert(model.Get{{  .Name  }}TableName()).
        		Columns(model.Get{{  .Name  }}FieldStringSlice()...).
        		Values(model.Get{{  .Name  }}AddField({{$short}})...).
        		PlaceholderFormat(d.PlaceHolder()).
        		ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
    if err != nil {
    	return
    }

    // run query
    log.Debug(ctx, "Add{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// BatchAdd{{ .Name }} batch add {{ .Name }} to the database
func (d *dao) BatchAdd{{ .Name }}s(ctx context.Context, list []*model.{{.Name}}) (err error) {
    if 0 == len(list) {
        return
    }

    if 50 < len(list) {
        return errors.New("batch operation supports up to 50")
    }

    // sql insert query, primary key provided by autoincrement
	sqlBuilder := squirrel.
        Insert(model.Get{{  .Name  }}TableName()).
        Columns(model.Get{{  .Name  }}FieldStringSlice()...)

    for _, {{ $short }} := range list {
        sqlBuilder = sqlBuilder.Values(model.Get{{  .Name  }}AddField({{$short}})...)
    }

    sqlStr, args, err := sqlBuilder.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return err
    }

     // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

    // run query
    log.Debug(ctx, "BatchAdd{{ .Name }}s", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Edit{{ .Name }} edit {{ .Name }} in the database
func (d *dao) Edit{{ .Name }}(ctx context.Context, setMap map[string]interface{}, condition map[string]interface{}) (err error) {
    // init build
    build := squirrel.Update(model.Get{{  .Name  }}TableName()).SetMap(setMap)

    // build sql
    build, err = d.AnalyticUpdate(build, condition)
    if err != nil {
       return
    }

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
       return
    }

     // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

	// run query
	log.Debug(ctx, "Edit{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
	if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Delete{{ .Name }} delete {{ .Name }} from the database
func (d *dao)  Delete{{ .Name }}(ctx context.Context,condition map[string]interface{}) (err error) {
    // init build
    build := squirrel.Delete(model.Get{{  .Name  }}TableName())

    // build sql
    build, err = d.AnalyticDelete(build, condition)
    if err != nil {
       return
    }

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
       return
    }

     // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

    // run query
    log.Debug(ctx, "Delete{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

/******************Generate By XO: read method segment**********************/
// Find{{ .Name }} query one {{ .Name }} by condition from the database
func (d *dao)  Find{{ .Name }}(ctx context.Context, condition map[string]interface{}) (res *model.{{ .Name }}, err error) {
    // init build
    build := squirrel.Select(model.Get{{ .Name }}FieldString()).From(model.Get{{  .Name  }}TableName())

    // forced master database query
    force := false
    if _, ok := condition["_force"]; ok {
        force = true
        delete(condition, "_force")
    }

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }
    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

    // run query
    log.Debug(ctx, "Find{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    res = &model.{{ .Name }}{}
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        if force {
            err = d.db.Master().QueryRow(ctx, sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
        } else {
            err = d.db.QueryRow(ctx, sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
        }
    } else {
        err = tx.QueryRow(sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
    }
    return
}

// Find{{ .Name }}s all {{ .Name }} by condition from the database
func (d *dao)  Find{{ .Name }}s(ctx context.Context, condition map[string]interface{}) (res []*model.{{ .Name }}, err error) {
    // init build
    build := squirrel.Select(model.Get{{ .Name }}FieldString()).From(model.Get{{  .Name  }}TableName())

    // forced master database query
    force := false
    if _, ok := condition["_force"]; ok {
        force = true
        delete(condition, "_force")
    }

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }
    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

    // run query
    log.Debug(ctx, "Find{{ .Name }}s", log.String("sql", fmt.Sprint(sqlStr, args)))
    var rows *sql.Rows
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        if force {
            rows, err = d.db.Master().Query(ctx, sqlStr, args...)
        } else {
            rows, err = d.db.Query(ctx, sqlStr, args...)
        }
    } else {
        rows, err = tx.Query(sqlStr, args...)
    }

    if err != nil {
        return
    }
    defer rows.Close()

    // load results
    res = make([]*model.{{ .Name }}, 0)
    for rows.Next() {
        // scan rows
        {{ $short }} := &model.{{ .Name }}{}
        err = rows.Scan(model.Get{{ .Name }}ScanField({{ $short }})...)
        if err != nil {
            return
        }

        res = append(res, {{ $short }})
    }
    return
}

// Count{{ .Name }} return count by condition from the database
func (d *dao)  Count{{ .Name }}(ctx context.Context, condition map[string]interface{}) (num int, err error) {
    // init build
    build := squirrel.Select("count(*)").From(model.Get{{  .Name  }}TableName())

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }

    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

     // parse sql to adapter databases
     sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
     if err != nil {
        return
     }

    // run query
    log.Debug(ctx, "Count{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        err = d.db.QueryRow(ctx, sqlStr, args...).Scan(&num)
    } else {
        err = tx.QueryRow(sqlStr, args...).Scan(&num)
    }
    return
}

// BatchCount{{ .Name }} groups records by groupKey and returns the count of each group.
// The result is a slice of maps, where each map's key is the groupKey field value
// and the value is the number of records in that group.
// Example: groupKey="status" returns [{1: 10}, {2: 5}] meaning 10 records with status=1, 5 with status=2.
func (d *dao)  BatchCount{{ .Name }}(ctx context.Context, groupKey string, condition map[string]interface{}) (res []map[string]int, err error) {
    // init build
    build := squirrel.Select("count(*) as num, " + groupKey).From(model.Get{{  .Name  }}TableName())

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }

    sqlStr, args, err := build.GroupBy(groupKey).PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = parser.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "BatchCount{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
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
    res = make([]map[string]int, 0)
    for rows.Next() {
    	// scan rows
    	var num int
    	var groupValue string
    	err = rows.Scan(&num, &groupValue)
    	if err != nil {
    		return
    	}

    	res = append(res, map[string]int{
    		groupValue: num,
    	})
    }
    return
}
