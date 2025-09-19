// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server for gin swagger example.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
package main

import (
	"go_backend/controllers"
	"go_backend/database"
	"go_backend/models"
	"go_backend/router"

	_ "go_backend/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"

	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @Summary 测试API
// @Description 测试API接口，返回用户列表
// @Tags 测试接口
// @Accept json
// @Produce json
// @Success 200 {array} models.User "用户列表"
// @Router /ping [get]
func testApi(ctx *gin.Context) {
	users := []models.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	ctx.JSON(200, users)
}

func main() {
	database.InitDB()
	db := database.GetDB()
	database.InitializeDefaultData(db)
	r := router.SetupRouter()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api")
	{
		api.GET("/ping", testApi)
		api.GET("/pagination", controllers.GetPagination)                 // 设备分页
		api.GET("/devices", controllers.GetDevices)                       // 获取设备
		api.GET("/deviceByid", controllers.GetDeviceByID)                 // 根据id获取设备
		api.POST("/update-device-status", controllers.UpdateDeviceStatus) // 根据设备id修改设备状态
		api.POST("/create-device", controllers.CreateDevice)              // 新增设备
	}
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
