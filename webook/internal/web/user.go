package web

import (
	"net/http"
	"time"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	birthdayExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		birthdayRegexPattern = `^(19|20)\d\d[- /.](0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	birthdayExp := regexp.MustCompile(birthdayRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		birthdayExp: birthdayExp,
	}
}

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 在这里登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 我可以随便设置值了
	// 你要放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type editReq struct {
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	var req editReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "解析编辑信息出错")
		return
	}
	if len(req.Nickname) == 0 || len(req.Nickname) > 10 {
		ctx.String(http.StatusOK, "昵称长度非法")
		return
	}
	if len(req.Introduction) == 0 || len(req.Introduction) > 100 {
		ctx.String(http.StatusOK, "简介长度非法")
		return
	}
	ok, err := u.birthdayExp.MatchString(req.Birthday)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式错误")
		return
	}

	sess := sessions.Default(ctx)
	id, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "获取登录id出错")
		return
	}
	if err := u.svc.Edit(ctx.Request.Context(), domain.User{
		Id:           id,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	}); err != nil {
		ctx.String(http.StatusOK, "更新信息出错")
	}
	ctx.String(http.StatusOK, "更新信息成功")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "获取登录id出错")
		return
	}
	user, err := u.svc.Profile(ctx.Request.Context(), id)
	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}
	type profileResp struct {
		Id           int64  `json:"id"`
		Email        string `json:"email"`
		Ctime        string `json:"create_time"`
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	ctx.JSON(http.StatusOK, profileResp{
		Id:           user.Id,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Birthday:     user.Birthday,
		Introduction: user.Introduction,
		Ctime:        time.UnixMilli(user.Ctime).Local().Format("2006-01-02 15:04:05"),
	})
}
