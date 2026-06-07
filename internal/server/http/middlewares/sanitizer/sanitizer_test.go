package sanitizer

// sanitizer_test.go
// 全场景测试：验证 Sanitizer 对各种请求结构的 XSS 清洗正确性
//
// 测试矩阵：
//   1.  简单 struct（普通字符串字段）
//   2.  纯数字字符串不被清洗
//   3.  空字符串不被清洗
//   4.  nil 指针字段不 panic
//   5.  *string 指针字段清洗
//   6.  嵌套 struct 递归清洗
//   7.  []string 切片每个元素清洗
//   8.  []struct 切片中结构体字段清洗
//   9.  []*struct 指针切片清洗
//  10.  map[string]interface{} 中字符串值清洗（BUG 2 修复验证）
//  11.  map[string]string 中字符串值清洗
//  12.  map[string]*struct 中结构体字段清洗
//  13.  复合嵌套：struct + []*struct + []string + map 同时存在
//  14.  nil 指针元素在 slice 中不 panic
//  15.  skipUrls 白名单路径跳过清洗
//  16.  安全 HTML 标签（<b>/<i>/<a>）被 UGC policy 保留
//  17.  危险标签（<script>/<iframe>/<img onerror>）被清除
//  18.  多层嵌套 struct 递归清洗

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// 辅助函数
// ─────────────────────────────────────────────

// newTestCtx creates a minimal *gin.Context for use in unit tests.
func newTestCtx(path string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	return ctx
}

// ─────────────────────────────────────────────
// 测试用 model
// ─────────────────────────────────────────────

type simpleReq struct {
	Name    string
	Age     string  // 纯数字，不应被清洗
	Bio     string
	Empty   string
	NilPtr  *string
}

type ptrReq struct {
	Title *string
}

type addressStruct struct {
	Street string
	City   string
}

type nestedReq struct {
	UserName string
	Addr     addressStruct
}

type deepNestedReq struct {
	Level1 string
	Child  nestedReq
}

type sliceStringReq struct {
	Tags []string
}

type itemStruct struct {
	Content string
	Count   string // 纯数字
}

type sliceStructReq struct {
	Items []itemStruct
}

type ptrSliceReq struct {
	Items []*itemStruct
}

type mapInterfaceReq struct {
	Params map[string]interface{}
}

type mapStringReq struct {
	Labels map[string]string
}

type mapPtrStructReq struct {
	UserMap map[string]*itemStruct
}

type complexChildStruct struct {
	Title  string
	Tags   []string
	Extra  map[string]interface{}
}

type complexReq struct {
	ID       string
	Children []*complexChildStruct
	Labels   map[string]string
}

// ─────────────────────────────────────────────
// 1. 简单 struct
// ─────────────────────────────────────────────

// TestSanitize_SimpleFields 验证普通字符串字段的 XSS 清洗。
func TestSanitize_SimpleFields(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &simpleReq{
		Name:  `<script>alert('xss')</script>johnsun`,
		Age:   "123",
		Bio:   `<b>hello</b><script>evil()</script>`,
		Empty: "",
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "johnsun", req.Name, "script tag should be stripped")
	assert.Equal(t, "123", req.Age, "numeric-only string should not be changed")
	assert.Equal(t, "<b>hello</b>", req.Bio, "safe HTML should be kept, unsafe removed")
	assert.Equal(t, "", req.Empty, "empty string should remain empty")
}

// ─────────────────────────────────────────────
// 2. 纯数字字符串不被清洗
// ─────────────────────────────────────────────

// TestSanitize_NumericString 验证纯数字字符串不被 sanitize 处理。
func TestSanitize_NumericString(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &simpleReq{Age: "9876543210"}
	s.Sanitize(ctx, req)
	assert.Equal(t, "9876543210", req.Age)
}

// ─────────────────────────────────────────────
// 3. 空字符串不被清洗
// ─────────────────────────────────────────────

// TestSanitize_EmptyString 验证空字符串不被处理。
func TestSanitize_EmptyString(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &simpleReq{Name: ""}
	s.Sanitize(ctx, req)
	assert.Equal(t, "", req.Name)
}

// ─────────────────────────────────────────────
// 4. nil 指针字段不 panic（BUG 1 修复验证）
// ─────────────────────────────────────────────

// TestSanitize_NilPointerField 验证 nil 指针字段不会导致 panic。
func TestSanitize_NilPointerField(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &simpleReq{Name: "ok", NilPtr: nil}
	assert.NotPanics(t, func() {
		s.Sanitize(ctx, req)
	}, "nil pointer field must not panic")
	assert.Equal(t, "ok", req.Name)
}

// ─────────────────────────────────────────────
// 5. *string 指针字段清洗
// ─────────────────────────────────────────────

// TestSanitize_PointerToString 验证 *string 字段的 XSS 清洗。
// bluemonday UGC policy 保留 <img> 标签但移除危险属性（如 onerror）。
func TestSanitize_PointerToString(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	title := `<img src=x onerror=alert(1)>hello`
	req := &ptrReq{Title: &title}
	s.Sanitize(ctx, req)

	// UGC policy keeps <img src="x"> but strips onerror attribute
	assert.Equal(t, `<img src="x">hello`, *req.Title, "onerror attribute should be stripped, img tag kept")
}

// ─────────────────────────────────────────────
// 6. 嵌套 struct 递归清洗
// ─────────────────────────────────────────────

// TestSanitize_NestedStruct 验证嵌套结构体中字段的递归清洗。
// bluemonday UGC policy 对 javascript: href 的处理：完全移除 <a> 标签，只保留文本内容。
func TestSanitize_NestedStruct(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &nestedReq{
		UserName: `<script>x</script>john`,
		Addr: addressStruct{
			Street: `<a href="javascript:evil()">click</a>`,
			City:   "Beijing",
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "john", req.UserName)
	// UGC policy strips the entire <a> tag when href is javascript:, leaving only text
	assert.Equal(t, "click", req.Addr.Street)
	assert.Equal(t, "Beijing", req.Addr.City)
}

// ─────────────────────────────────────────────
// 7. []string 切片每个元素清洗
// ─────────────────────────────────────────────

// TestSanitize_SliceOfStrings 验证字符串切片中每个元素都被清洗。
func TestSanitize_SliceOfStrings(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &sliceStringReq{
		Tags: []string{
			"golang",
			`<script>evil()</script>`,
			"<b>bold</b>",
			"",
			"123",
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "golang", req.Tags[0])
	assert.Equal(t, "", req.Tags[1], "script tag should be stripped entirely")
	assert.Equal(t, "<b>bold</b>", req.Tags[2], "safe bold tag should be kept")
	assert.Equal(t, "", req.Tags[3], "empty string should remain empty")
	assert.Equal(t, "123", req.Tags[4], "numeric string should not change")
}

// ─────────────────────────────────────────────
// 8. []struct 切片中结构体字段清洗
// ─────────────────────────────────────────────

// TestSanitize_SliceOfStructs 验证结构体切片中每个元素的字段都被清洗。
func TestSanitize_SliceOfStructs(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &sliceStructReq{
		Items: []itemStruct{
			{Content: "hello", Count: "10"},
			{Content: `<iframe src="evil.com"></iframe>`, Count: "20"},
			{Content: `<b>bold</b>`, Count: "30"},
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "hello", req.Items[0].Content)
	assert.Equal(t, "10", req.Items[0].Count, "numeric Count should not change")
	assert.Equal(t, "", req.Items[1].Content, "iframe should be stripped by UGC policy")
	assert.Equal(t, "<b>bold</b>", req.Items[2].Content)
}

// ─────────────────────────────────────────────
// 9. []*struct 指针切片清洗
// ─────────────────────────────────────────────

// TestSanitize_SliceOfPtrStructs 验证指针结构体切片中每个元素的字段都被清洗。
func TestSanitize_SliceOfPtrStructs(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &ptrSliceReq{
		Items: []*itemStruct{
			{Content: `<script>xss()</script>clean`},
			{Content: "<b>safe</b>"},
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "clean", req.Items[0].Content)
	assert.Equal(t, "<b>safe</b>", req.Items[1].Content)
}

// ─────────────────────────────────────────────
// 10. map[string]interface{} 中字符串值清洗（BUG 2 修复验证）
// ─────────────────────────────────────────────

// TestSanitize_MapStringInterface 验证 map[string]interface{} 中字符串值被清洗。
func TestSanitize_MapStringInterface(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &mapInterfaceReq{
		Params: map[string]interface{}{
			"name":    `<script>xss()</script>johnsun`,
			"count":   "42",
			"safe":    "<b>bold</b>",
			"empty":   "",
			"numeric": "100",
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "johnsun", req.Params["name"], "script tag should be stripped from map value")
	assert.Equal(t, "42", req.Params["count"], "numeric string in map should not be changed")
	assert.Equal(t, "<b>bold</b>", req.Params["safe"], "safe HTML in map should be kept")
	assert.Equal(t, "", req.Params["empty"], "empty string in map should remain empty")
	assert.Equal(t, "100", req.Params["numeric"])
}

// ─────────────────────────────────────────────
// 11. map[string]string 中字符串值清洗
// ─────────────────────────────────────────────

// TestSanitize_MapStringString 验证 map[string]string 中字符串值被清洗。
func TestSanitize_MapStringString(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &mapStringReq{
		Labels: map[string]string{
			"key1": `<script>alert(1)</script>clean`,
			"key2": "normal",
			"key3": "456",
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "clean", req.Labels["key1"])
	assert.Equal(t, "normal", req.Labels["key2"])
	assert.Equal(t, "456", req.Labels["key3"], "numeric map value should not change")
}

// ─────────────────────────────────────────────
// 12. map[string]*struct 中结构体字段清洗
// ─────────────────────────────────────────────

// TestSanitize_MapPtrStruct 验证 map value 为 *struct 时，结构体字段被递归清洗。
func TestSanitize_MapPtrStruct(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &mapPtrStructReq{
		UserMap: map[string]*itemStruct{
			"u1": {Content: `<script>evil()</script>alice`, Count: "10"},
			"u2": {Content: "<b>bob</b>", Count: "20"},
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "alice", req.UserMap["u1"].Content)
	assert.Equal(t, "10", req.UserMap["u1"].Count, "numeric Count should not change")
	assert.Equal(t, "<b>bob</b>", req.UserMap["u2"].Content)
}

// ─────────────────────────────────────────────
// 13. 复合嵌套：struct + []*struct + []string + map 同时存在
// ─────────────────────────────────────────────

// TestSanitize_ComplexNested 验证复杂嵌套结构的全面清洗。
func TestSanitize_ComplexNested(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &complexReq{
		ID: "999", // 纯数字，不应变化
		Children: []*complexChildStruct{
			{
				Title: `<script>evil()</script>child1`,
				Tags:  []string{"tag1", `<img onerror=alert(1) src=x>`},
				Extra: map[string]interface{}{
					"note":    `<b>important</b><script>bad()</script>`,
					"numeric": "100",
				},
			},
		},
		Labels: map[string]string{
			"env": `<script>hack()</script>prod`,
		},
	}

	assert.NotPanics(t, func() {
		s.Sanitize(ctx, req)
	})

	assert.Equal(t, "999", req.ID, "numeric ID should not change")
	assert.Equal(t, "child1", req.Children[0].Title, "script in nested struct should be stripped")
	assert.Equal(t, "tag1", req.Children[0].Tags[0])
	// UGC policy keeps <img> but strips onerror attribute
	assert.Equal(t, `<img src="x">`, req.Children[0].Tags[1], "onerror should be stripped, img tag kept")
	assert.Equal(t, "<b>important</b>", req.Children[0].Extra["note"], "script in map value should be stripped")
	assert.Equal(t, "100", req.Children[0].Extra["numeric"], "numeric map value should not change")
	assert.Equal(t, "prod", req.Labels["env"], "script in map[string]string should be stripped")
}

// ─────────────────────────────────────────────
// 14. nil 指针元素在 slice 中不 panic
// ─────────────────────────────────────────────

// TestSanitize_NilElementInSlice 验证 slice 中的 nil 指针元素不会导致 panic。
func TestSanitize_NilElementInSlice(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &ptrSliceReq{
		Items: []*itemStruct{
			{Content: "ok"},
			nil, // nil 指针元素
			{Content: `<script>xss</script>clean`},
		},
	}

	assert.NotPanics(t, func() {
		s.Sanitize(ctx, req)
	}, "nil element in slice must not panic")

	assert.Equal(t, "ok", req.Items[0].Content)
	assert.Nil(t, req.Items[1])
	assert.Equal(t, "clean", req.Items[2].Content)
}

// ─────────────────────────────────────────────
// 15. skipUrls 白名单路径跳过清洗
// ─────────────────────────────────────────────

// TestSanitize_SkipUrl 验证 skipUrls 中的路径不会被清洗。
func TestSanitize_SkipUrl(t *testing.T) {
	original := skipUrls
	skipUrls = []string{"/api/skip"}
	defer func() { skipUrls = original }()

	s := NewSanitizer()
	ctx := newTestCtx("/api/skip")

	req := &simpleReq{
		Name: `<script>alert(1)</script>`,
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, `<script>alert(1)</script>`, req.Name,
		"skipped URL should not sanitize the request")
}

// ─────────────────────────────────────────────
// 16. 安全 HTML 标签被 UGC policy 保留
// ─────────────────────────────────────────────

// TestSanitize_SafeHTMLKept 验证 UGC policy 允许的安全 HTML 标签被保留。
// 注意：bluemonday UGC policy 会为外部链接自动添加 rel="nofollow"。
func TestSanitize_SafeHTMLKept(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &simpleReq{
		Name: "<b>bold</b> and <i>italic</i>",
		Bio:  `<a href="https://example.com">link</a>`,
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "<b>bold</b> and <i>italic</i>", req.Name, "b and i tags should be kept")
	// UGC policy adds rel="nofollow" to external links
	assert.Equal(t, `<a href="https://example.com" rel="nofollow">link</a>`, req.Bio, "safe href kept with rel=nofollow added")
}

// ─────────────────────────────────────────────
// 17. 危险标签被清除
// ─────────────────────────────────────────────

// TestSanitize_DangerousTagsStripped 验证各种危险 XSS 载体被清除。
// 断言基于 bluemonday UGC policy 的实际行为：
//   - <img> 标签被保留，但危险属性（onerror）被移除
//   - <a href="javascript:"> 整个标签被移除，只保留文本
//   - <script>/<iframe>/<svg onload>/<object> 完全被移除
func TestSanitize_DangerousTagsStripped(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	cases := []struct {
		input    string
		expected string
		desc     string
	}{
		{`<script>alert(1)</script>`, "", "script tag stripped entirely"},
		{`<iframe src="evil.com"></iframe>`, "", "iframe tag stripped entirely"},
		{`<img src=x onerror=alert(1)>`, `<img src="x">`, "img kept but onerror attribute stripped"},
		{`<svg onload=alert(1)>`, "", "svg with onload stripped entirely"},
		{`<a href="javascript:evil()">click</a>`, "click", "javascript: href causes entire <a> to be stripped"},
		{`<object data="evil.swf"></object>`, "", "object tag stripped entirely"},
	}

	for _, tc := range cases {
		req := &simpleReq{Name: tc.input}
		s.Sanitize(ctx, req)
		assert.Equal(t, tc.expected, req.Name, tc.desc)
	}
}

// ─────────────────────────────────────────────
// 18. 多层嵌套 struct 递归清洗
// ─────────────────────────────────────────────

// TestSanitize_DeepNestedStruct 验证多层嵌套结构体的递归清洗。
func TestSanitize_DeepNestedStruct(t *testing.T) {
	s := NewSanitizer()
	ctx := newTestCtx("/api/test")

	req := &deepNestedReq{
		Level1: `<script>l1</script>level1`,
		Child: nestedReq{
			UserName: `<script>l2</script>level2`,
			Addr: addressStruct{
				Street: `<script>l3</script>street`,
				City:   "Shanghai",
			},
		},
	}
	s.Sanitize(ctx, req)

	assert.Equal(t, "level1", req.Level1)
	assert.Equal(t, "level2", req.Child.UserName)
	assert.Equal(t, "street", req.Child.Addr.Street)
	assert.Equal(t, "Shanghai", req.Child.Addr.City)
}
