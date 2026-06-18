package controllers

import (
	"net/http"
	"shenpi/models"
	"shenpi/services"
	"shenpi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(models.DB),
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
		Role:     "user",
	}

	if err := c.userService.Register(user); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "注册失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.Success(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (c *UserController) Profile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	user, err := c.userService.GetByID(userID)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(ctx, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (c *UserController) List(ctx *gin.Context) {
	keyword := ctx.Query("keyword")

	users, err := c.userService.GetList(keyword)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "获取用户列表失败: "+err.Error())
		return
	}

	result := make([]gin.H, 0, len(users))
	for _, user := range users {
		result = append(result, gin.H{
			"id":   user.ID,
			"name": user.Name,
			"username": user.Username,
		})
	}

	utils.Success(ctx, result)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的ID")
		return
	}

	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(ctx, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"role":     user.Role,
	})
}
