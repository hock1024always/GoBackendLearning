# Gin

## 使用Group进行版本选择

```go
func initCourse(r *gin.Engine) {
	v1 := r.Group("/v1", middleware.TokenCheck, middleware.AuthCheck)
	v1.POST("/course", web.CreateCourse)
	v1.GET("/course", web.GetCourse)
	v1.PUT("/course", web.EditCourse)
	v1.DELETE("/course", web.DeleteCourse)

	v2 := r.Group("/v2")
	v2.POST("/course", web.CreateCourseV2)
	v2.PUT("/course", web.EditCourseV2)
}
```

## 参数校验

| 标签      | 用途              | 示例                 |
| :-------- | :---------------- | :------------------- |
| `form`    | 绑定表单/查询参数 | `form:"user_name"`   |
| `json`    | 绑定 JSON 数据    | `json:"userName"`    |
| `xml`     | 绑定 XML 数据     | `xml:"UserName"`     |
| `binding` | 定义验证规则      | `binding:"required"` |

```go
type registerReq struct {
	UserName string `form:"user_name" binding:"required"`    //必填
	Password string `form:"pwd" binding:"required"`          //必填
	Phone    string `form:"phone" binding:"required e164"`   //手机号格式
	Email    string `form:"email" binding:"omitempty,email"` //不必填，如果为空忽略
}

func Register(c *gin.Context) {
	req := &registerReq{}
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, req)
}
```

