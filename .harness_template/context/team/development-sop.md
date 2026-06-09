# 团队开发标准 SOP (Standard Operating Procedure)

> **核心原则**：各服务代码实现必须遵循本标准套路，严禁 AI 或开发者在没有达成共识的情况下自行发挥。
>
> **特别强调：绝不瞎篡改、自创代码规范，具体服务的架构约束完全以其目录下的 `README.md` 为准！**

## 1. 架构分层与实现路径

业务功能的实现必须在设计门禁时确认，编码时严格遵守以下各层级的规范套路（以 `api/{service-name}` 架构为范本）：

### 1.1 数据访问层 (Model / IFace / DAO)
底层数据访问层采用接口隔离与依赖注入的设计，**必须通过 `xo` 等组件自动生成**：
- **Model 层 (`internal/model`)**：只存放与数据库表一一对应的实体结构体，带有对应的 json 标签。**注意：Model 字段类型由 DDL 字段类型决定，Go 代码中会出现 uint / uint64 等 unsigned 类型，绝对禁止因为业务层的出入参定义而篡改 model 的定义。**
- **IFace 层 (`internal/dao/iface`)**：定义 DAO 层必须实现的底层操作接口，用于实现接口隔离。
- **DAO 注册与组合 (`internal/dao/dao.go`)**：在基础 DAO 文件中，必须将 `iface` 定义的各个领域接口组合进主 `Dao` 接口中，以便 Service 层统一注入使用。
  ```go
  // 防错示例（注意只在 Dao interface 中进行组合，而不是 struct 中）：
  type Dao interface {
      Close() error
      Ping(ctx context.Context) error
      // ... 基础组件方法
      Begin(ctx context.Context) (context.Context, error)

      // 组合 iface 层的抽象接口 (在此处增加你新生成的接口)
      iface.TUser
      // iface.YourNewModel 
  }
  ```
- **DAO 实现层 (`internal/dao`)**：负责具体的 SQL 拼接（使用 `squirrel` 库）与数据库交互。**注意：Add新增数据不能定义返回结果获取ID，比如从 `result.LastInsertId()` 获取ID，因为有些国产数据库并不支持返回插入的ID。**所有 SQL 操作必须支持基于上下文传递的事务传递能力（检查 `tx` 的存在）。

> ⚠️ **DAO condition 常量红线**：
> Service 层调用 DAO 时，condition map 中的特殊操作 key 和排序字段名**禁止硬编码字符串**，必须使用常量：
>
> ```go
> // ❌ 错误写法：硬编码字符串
> s.dao.FindTGroups(ctx, map[string]interface{}{
>     "_orderBy": "sort_order asc",
>     "_limit":   uint64(1),
> })
>
> // ✅ 正确写法：使用常量
> s.dao.FindTGroups(ctx, map[string]interface{}{
>     drivers.OpAction__order_by.String(): model.GetTGroupColumns().SortOrder + " asc",
>     drivers.OpAction__limit.String():    uint64(1),
> })
> ```
>
> 常用 drivers 操作常量：`OpAction__order_by`、`OpAction__limit`、`OpAction__offset`、`OpAction__group_by`、`OpAction__like`
> 字段名一律从 `model.GetXxxColumns().FieldName` 获取

> **AI 代码生成强制约束**：
> 1. 凡是涉及 DDL 或 DML 变更导致需要更新 DAO、IFace 和 Model 代码时，**强烈建议**通过 XO 命令行工具在本地生成。
> 2. 在缺少 `xo` 工具自动生成时，AI 若被要求手写 DAO/Model/IFace 或 Ecode，**必须逐步加载以下默认模板文件，并严格遵循其写法示例进行 1:1 像素级复刻**。无论是结构、方法签名还是逻辑组织，必须与 demo 保持绝对一致，严禁自我创造！**（唯一例外：项目内部的包导入路径必须智能替换为当前工作目录的实际地址，但 `github.com` 等第三方包必须原样保留，避免编译失败）**
>    - DAO IFACE 模板: `api/{service-name}/internal/dao/iface/tuser.xo.go`
>    - DAO 实现模板: `api/{service-name}/internal/dao/tuser.xo.go`
>    - MODEL 模板: `api/{service-name}/internal/model/tuser.xo.go`
>    - 错误码(ecode)模板: `api/{service-name}/internal/ecode/i18n.go`
>    *(注：`{service-name}` 为当前工作的服务名，任一初始化的服务均会默认包含这些 user 相关的 demo 文件，应就近读取。)*
> 3. AI 在手写补充 DAO 方法（仅限非 XO 自动生成的定制逻辑）时，必须同步维护 `iface` 的接口定义并在 `dao` 目录下实现，保持模式的一致性。

### 1.2 gRPC 接口服务 (gRPC Service)
- **场景**：服务间的高效内部 RPC 调用。
- **套路规范**：
  - **IDL 定义与生成**：在 `api/idl/{service_name}/` 找相应的 proto 文件定义接口，并**必须主动执行** `waterdrop protoc --grpc --swagger *.proto` 生成新的 stub 文件。
  - **服务注册**：在 `internal/server/grpc/server.go` 中将实现类注册到 gRPC Server。
  - **业务实现 (`internal/service`)**：在 `service` 目录下实现 proto 中定义的 RPC 接口。
    - **核心套路**：**开启事务(按需) -> defer 捕获错误并统一 Rollback -> 核心业务逻辑及模型组装 -> 调用 DAO -> 所有操作成功后显式 Commit -> 组装 Reply**。

> 📄 **gRPC Service 及事务处理实现样板代码**

```go
package service

import (
	"coder/api/idl/user"
	"coder/api/user/internal/model"
	"context"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/UnderTreeTech/waterdrop/pkg/utils/ecode"
	"github.com/UnderTreeTech/waterdrop/pkg/utils/xtime"
)

// GetUserInfo 获取用户信息业务逻辑样板：包含入参校验、调用 DAO 层进行查询并转换错误码
func (s *Service) GetUserInfo(ctx context.Context, req *user.UserInfoReq) (reply *user.UserInfoReply, err error) {
	reply = &user.UserInfoReply{}
	reply.UserInfo = &user.UserInfo{}

	// 构建查询条件
	cond := make(map[string]interface{})
	cond[model.GetTUserColumns().UserID] = req.GetUid()
	
	// 调用 DAO 层执行查询
	userModel, err := s.dao.FindTUser(ctx, cond)
	if uerr != nil {
		// 打印日志并转换对外错误码
		log.Error(ctx, "get user info fail", log.Uint64("uid", req.GetUid()),
			log.String("error", err.Error()))
		return
	}

	// 组装并赋值 Reply 数据 (此处仅作示例)
	reply.UserInfo.Uid = userModel.UserID
	reply.UserInfo.UserName = userModel.UserName

	return reply, nil
}

// AddUserBiz 添加用户业务逻辑样板：处理复杂的写操作，涉及多表关联写入与事务处理
func (s *Service) AddUserBiz(ctx context.Context, req *user.UserInfoBizReq) (reply *user.BatchGetUidsReply, err error) {
	reply = &user.BatchGetUidsReply{
		Uids: make([]uint64, 0),
	}

	// 开启事务
	if ctx, err = s.dao.Begin(ctx); err != nil {
		log.Error(ctx, "add user fail", log.Any("ctx", ctx),
			log.String("error", err.Error()))
		return
	}

	// 新增用户属于复杂写操作，开启事务。
	// 注意：使用 defer 统一处理回滚，而 Commit 必须在所有操作完成的最后显式调用，保证逻辑清晰通顺。
	defer func() {
		if err != nil {
			log.Error(ctx, "add user fail, rollbacked", log.String("error", err.Error()))
			if rerr := s.dao.Rollback(ctx); rerr != nil {
				log.Error(ctx, "rollback fail", log.Any("ctx", ctx), log.String("error", rerr.Error()))
			}
		}
	}()

	userInfo := &model.TUser{
		UserID:      uint64(xtime.Now().CurrentUnixTime()),
		UserName:    req.Name,
		UserAccount: req.Account,
	}

	// 调用 DAO 层插入数据
	err = s.dao.AddTUser(ctx, userInfo)
	if err != nil { 
	    return // err 非空，触发 defer Rollback 
	}
	
	// ... 省略其余关联写入 ...

	// 所有表操作均已成功结束，最后显式提交事务
	if err = s.dao.Commit(ctx); err != nil {
		log.Error(ctx, "commit fail", log.String("error", err.Error()))
		return
	}

	return
}
```
### 1.3 HTTP 接口服务 (HTTP Controller)
- **场景**：直接暴露给前端、第三方或内部基于 HTTP 的通信。
- **套路规范**：HTTP Service 强调**内外隔离**，Controller 绝对禁止出现 SQL 或者直接组装 DB Model，所有输入输出必须定义专门的 Model。
  - **1. 路由注册**：在 `internal/server/http/router.go` 中注册 HTTP 路由。
  - **2. 定义入参/出参 (`http/model`)**：在 `internal/server/http/model/` 目录下（或指定的 http dto 目录），专门定义该接口的 Request 和 Response 结构体。**注意：HTTP出入参model定义仅允许创建在internal/server/http/model/目录下**。
  - **3. Controller 实现 (`http` 目录)**：在 `internal/server/http/{module}.go` 中实现控制器逻辑（如 `http/user.go`）。只负责：`解析 Gin 上下文参数 -> 转换为业务参数 -> 调用 internal/service 的逻辑层方法 -> 封装标准响应返回`。**严禁写重度业务逻辑或直接调 DAO**。
  - **4. 业务下沉 (`internal/service`)**：必须将核心业务下沉到 `internal/service` 中的普通 Service 方法中处理。**注意：写service逻辑有需要定义struct使用时，必须放在 `{harness_name}/api/{service_name}/internal/model/` 目录下，禁止在业务service文件中定义struct。**业务service需要引用HTTP请求的出入参model定义时，重命名HTTP出入参定义model目录导入命名，如 `m "github.com/UnderTreeTech/layout/internal/server/http/model"`。service函数非复杂接口，入参尽量不要超过3个。

> ⚠️ **HTTP 路由风格红线**：禁止 RESTful 路径参数**：不得使用 `/api/group/:id`、`/api/agent/:id/move` 这类把资源 ID 放在 URL 路径中的风格。

> 📄 **HTTP Service 完整闭环样板代码**

**3.1 路由注册 (`internal/server/http/router.go`)**
```go
package http

import (
	"github.com/gin-gonic/gin"
)

// registerAPI 注册API
func registerAPI(engine *gin.Engine) {
	login := engine.Group("/api")
	{
		login.GET("/app/user", getUserInfo) // 获取用户信息
	}
}
```

**3.2 HTTP 专属出入参 (`internal/server/http/model/user.go`)**
用于与外部前端进行 JSON 交互，不污染底层数据库 Model。**注意：HTTP出入参model定义仅允许创建在internal/server/http/model/目录下**
```go
package model

// GetUserInfoReq 获取用户信息请求
type GetUserInfoReq struct {
	UserId string `json:"user_id" form:"user_id" validate:"required"`
	UserName string `json:"user_name" form:"user_name"`
}

// GetUserInfoReply 获取用户信息返回
type GetUserInfoReply struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}
```

**3.3 HTTP Controller (`internal/server/http/user.go`)**
只负责 HTTP 上下文处理：`获取并绑定参数 -> 调用 Service 的逻辑层方法 -> 通过统一 reply 格式化返回`。
```go
package http

import (
	"net/http"

	"coder/api/user/internal/server/http/model"
	"coder/api/user/internal/utils/reply"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// getUserInfo 查询用户信息
func getUserInfo(ctx *gin.Context) {
	req := &model.GetUserInfoReq{}
	// 绑定并校验 HTTP 请求参数
	if err := ShouldBind(ctx, &req); err != nil {
		return
	}

	// 调用具体的业务方法（业务逻辑下沉到 svc）
	// 注意：严禁在这里直接写业务核心判断，或直接调用 DAO
	resp, err := svc.GetUserInfo(ctx.Request.Context(), cast.ToUint64(req.UserId))
	
	// 封装标准返回结构
	ctx.JSON(http.StatusOK, reply.Reply(ctx, resp, err))
}
```

**3.4 service实现 (`internal/service/user.go`)**
处理真正的业务逻辑：`获取并绑定参数 -> 调用 Service 的逻辑层方法 -> 通过统一 reply 格式化返回`。
**注意：写service逻辑有需要定义struct使用时，必须放在 `{harness_name}/api/{service_name}/internal/model/` 目录下，禁止在业务service文件中定义struct。**
**注意：业务service需要引用HTTP请求的出入参model定义时，重命名HTTP出入参定义model目录导入命名，如 `m "github.com/UnderTreeTech/layout/internal/server/http/model"`。**
**注意：service函数非复杂接口，入参尽量不要超过3个**

```go
package service

import (
	"context"
	m "github.com/UnderTreeTech/layout/internal/server/http/model"
)

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, uid string) (reply *m.GetUserInfoReply, err error) {
	// 业务代码
	return &m.GetUserInfoReply{
		UserId:   "1",
		UserName: "johnsun",
	}, nil
}
```
### 1.4 跨服务 RPC 调用 (gRPC Client)
- **场景**：当业务逻辑需要进行跨服务调用（调用其他 gRPC Service）时，无论是 HTTP Service 还是 gRPC Service 均遵循以下规范。
- **规范**：所有外部服务的 RPC Client 在 `Service` 结构体中统一维护，在 `New` 初始化函数中通过 `waterdrop` 的客户端库创建连接，绑定必要的拦截器（如断路器），并最终挂载到 `Service` 实例。

> 📄 **RPC Client 定义与初始化样板代码**

所有外部服务的 RPC Client 在 `Service` 结构体中统一维护，在 `New` 初始化函数中通过 `waterdrop` 的客户端库创建连接，绑定必要的拦截器（如断路器），并最终挂载到 `Service` 实例。

```go
package service

import (
	"fmt"

	// 引入依赖外部服务的协议生成包
	"api/idl/user"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/client"
	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/config"
	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/interceptors"
)

type Service struct {
	// ...
	// 外部 RPC 客户端实例，按照依赖的服务隔离保存
	user user.UserClient
}

func New(d dao.Dao, cfg *Config) *Service {
	// 1. 初始化客户端配置并解析 yaml 配置
	userCfg := &config.ClientConfig{}
	if err := conf.Unmarshal("client.rpc.user", userCfg); err != nil {
		panic(fmt.Sprintf("unmarshal user client config fail, err msg %s", err.Error()))
	}

	// 2. 利用配置创建基础 RPC 客户端连接
	userRPCCli := client.New(userCfg)
	user := organization.NewUserClient(userRPCCli.GetConn())

	svc := &Service{
		// ...
		user: user,
	}

	return svc
}
```

### 4.2 RPC Client 业务调用 (`internal/service/group.go` 或其他业务文件)

在具体的业务逻辑方法中，直接通过 `Service` 挂载的特定客户端对象（例如 `s.user`）调用外部服务的 RPC 方法。

```go
	// 1. 构建跨服务 RPC 请求包
	userInfoReq := &user.UserInfoReq{
		Uid: req.GetUserId(),
	}

	// 2. 发起跨服务调用并传递 context (重要，为了链路追踪与超时控制)
	hostInfo, err := s.user.GetUserInfo(ctx, userInfoReq)
	if err != nil {
		log.Error(ctx, "get user info fail", log.Uint64("uid", groupInfo.HostUID),
			log.String("error", err.Error()))
		return
	}
```
### 1.5 内部普通 Service 方法
- **场景**：服务内部的公共业务逻辑抽取，供 HTTP Controller、gRPC Service 或定时任务、MQ 消费者调用。
- **套路规范**：直接在 `internal/service/` 下建立对应的业务文件实现即可，纯 Go 方法，无对外网络协议绑定。

> ⚠️ **Service 层错误变量命名红线**：
> - 函数使用 named return `err` 时，函数体内**禁止**使用 `ferr`/`qerr`/`cerr`/`aerr`/`uerr` 等别名。
> - Go 规范允许 `:=` 在函数体内**重声明** named return 参数（只要至少有一个新变量），因此直接用 `err` 即可：
>
> ```go
> // ❌ 错误写法：不必要的别名
> groups, ferr := s.dao.FindTGroups(ctx, cond)
> if ferr != nil { err = ecode.OperationFailed; return }
>
> // ✅ 正确写法：:= 合法重声明 named return err
> groups, err := s.dao.FindTGroups(ctx, cond)
> if err != nil { err = ecode.OperationFailed; return }
> ```
>
> **唯一例外**：defer 闭包中 `if err := s.dao.Rollback(ctx); err != nil {` 创建新作用域局部变量，不影响外层 named return，这是正确的。

## 2. 核心红线 (AI 必须遵守)

1. **以服务级文档为尊**：`api/{service_name}/README.md` 是该服务的绝对规范源，任何生成代码的结构、命名、分层均需照抄该文档的范式，**不得发明所谓的 biz 层或凭空改造目录结构**。
2. **规范 IDL 编写**：对于 gRPC 接口，AI 应主动基于 `context/team/protobuf-style-guide.md` 规范给出 proto 变更，必须提示执行 `protoc` 生成代码命令。
3. **大仓导入路径准确性**：在生成或修改 Go 导入路径时，必须匹配当前大仓 `go.mod` 中定义的实际 module 名称及当前工作目录地址。**注意：在进行 1:1 像素级复制模板代码时，必须将模板中的包导入地址（如 `coder/api/user/internal/...`）智能替换为当前工作目录的准确地址（注意：`github.com` 等第三方外部包必须原样保留，绝不允许替换），严禁生搬硬套导致编译失败。**
4. **全局错误码统一与服务级覆盖**：错误码统一定义在外部的 `api/ecode` 服务中，各服务必须引用该全局服务的错误码。**公共错误码段（100000-100999）已包含"内部错误"、"参数非法"等通用错误。**各服务内部的 `internal/ecode` 目录用来存放特有自定义错误码，或通过相同的 code 不同文案来实现对全局错误码文案的覆盖。**注意：在服务内的 `i18n.go` 中，必须优先注册全局错误码，然后再注册自定义错误码进行覆盖，详见 `context/team/error-code.md`。**
5. **DDL 字段类型必须对照规范**：编写 DDL 时必须逐字段对照 `context/team/db-design.md`，主键用 `bigint unsigned`，时间戳用 `bigint unsigned`，非负数用 `unsigned` 变体，所有字段 `NOT NULL` 有默认值，DDL 仅保留主键索引。
6. **DAO condition 禁止硬编码**：condition map 中特殊操作 key 必须用 `drivers.OpAction__xxx.String()` 常量，排序字段名必须从 `model.GetXxxColumns().FieldName` 获取，禁止硬编码 `"_orderBy"` 或 `"sort_order asc"` 等字符串。
7. **Service 层错误变量统一用 `err`**：函数体内 named return 参数 `err` 可通过 `:=` 合法重声明，禁止使用 `ferr`/`qerr`/`cerr` 等别名。唯一例外是 defer 闭包内的 `if err :=` 局部变量。

## 3. 代码收尾要求

1. **自动格式化**：在完成一个功能模块的代码编写后，必须主动执行 `go fmt ./...` 对当前修改的代码进行格式化。
2. **存根更新**：所有 `.proto` 的变更，必须在代码编写后、提交前，确保执行过生成命令，保证存根与定义一致。
