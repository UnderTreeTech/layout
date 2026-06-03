package parser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/format"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/valyala/bytebufferpool"
)

var (
	DatabaseTypeDM       = "dm"
	DatabaseTypeKingbase = "kingbase"
)

// ParserPool parser pool definition
type ParserPool struct {
	pool *sync.Pool
}

// NewParserPool creates a new ParserPool
func NewParserPool() *ParserPool {
	return &ParserPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return parser.New()
			},
		},
	}
}

// Get gets a parser from the ParserPool
func (pp *ParserPool) Get() *parser.Parser {
	p := pp.pool.Get().(*parser.Parser)
	return p
}

// Put returns the given parser to the ParserPool
// parser底层每次解析时，会重置parser的相关数据
func (pp *ParserPool) Put(buffer *parser.Parser) {
	pp.pool.Put(buffer)
}

var pp = NewParserPool()

const (
	QuestionFormat = "?"  // MYSQL系
	DollarFormat   = "$"  // PG系：PG、KINGBASE、GAUSE、CERDB
	ColonFormat    = ":"  // ORACLE系：ORACLE、DM
	AtpFormat      = "@p" // SQLSERVER

	BackQuotesFlags = format.RestoreNameBackQuotes | format.RestoreStringSingleQuotes | format.RestoreSpacesAroundBinaryOperation | format.RestoreKeyWordUppercase
	DoubleQuotes    = format.RestoreNameDoubleQuotes | format.RestoreStringSingleQuotes | format.RestoreSpacesAroundBinaryOperation | format.RestoreKeyWordUppercase
)

// QuoteSQL 处理SQL语句，目前仅DM数据库需要做语句特殊处理
func QuoteSQL(driver string, sql string) (quotedSQL string, err error) {
	switch driver {
	case DatabaseTypeDM:
		quotedSQL, err = quoteSQL(sql, DoubleQuotes)
		if err != nil {
			return
		}
		quotedSQL = replacePositionalPlaceholders(quotedSQL, ColonFormat)
	default:
		return sql, nil
	}
	return
}

// quoteSQL 字段名、表名加双引号或反引号，关键字大小写
func quoteSQL(oldStatement string, flags format.RestoreFlags) (newStatement string, err error) {
	p := pp.Get()
	defer pp.Put(p)

	newStatement = oldStatement
	rootNode, err := p.ParseOneStmt(newStatement, "", "")
	if err != nil {
		return
	}
	builder := &strings.Builder{}
	ctx := format.NewRestoreCtx(flags, builder)
	err = rootNode.Restore(ctx)
	if err != nil {
		return
	}
	builder.WriteByte(';')
	newStatement = builder.String()
	return
}

// replacePositionalPlaceholders 替换语句占位符
// 由于parser库只能解析?占位符，如适配其他数据库如DM、KINGBASE，生成的SQL的默认占位符最好都是?，不然quoteSQL解析会报错
func replacePositionalPlaceholders(sql, placeholder string) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			fmt.Fprintf(buf, "%s%d", placeholder, i)
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String()
}
