package controllers

import (
	"go_backend/database"
	"go_backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDevices godoc
// @Summary 获取所有设备列表
// @Description 获取系统中所有设备的信息列表
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {array} models.Device "设备列表"
// @Failure 500 {object} map[string]string "内部服务器错误"
// @Router /devices [get]
func GetDevices(c *gin.Context) {
	var devices []models.Device
	db := database.GetDB()

	// 查询所有设备
	if err := db.Find(&devices).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve devices"})
		return
	}

	// 返回设备列表作为 JSON 响应
	c.JSON(200, devices)
}

// GetDeviceByID godoc
// @Summary 根据ID获取设备信息
// @Description 根据设备ID获取单个设备的详细信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id query string true "设备ID"
// @Success 200 {object} map[string]interface{} "设备信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "设备未找到"
// @Router /device [get]
func GetDeviceByID(c *gin.Context) {
	// 获取 query 参数 id
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(400, gin.H{"error": "Missing device ID"})
		return
	}

	// 转换 id 为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid device ID"})
		return
	}

	db := database.GetDB()
	var device models.Device

	// 查询设备
	if err := db.First(&device, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	// 返回设备信息
	c.JSON(200, gin.H{
		"id":       device.ID,
		"name":     device.Name,
		"status":   device.Status,
		"position": device.Position,
	})
}

// UpdateDeviceStatus godoc
// @Summary 更新设备状态
// @Description 根据设备ID更新设备状态
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param request body map[string]int true "设备状态更新请求" SchemaExample({"id": 1, "status": 1})
// @Success 200 {object} map[string]interface{} "更新成功响应"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "设备未找到"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /device/status [put]
func UpdateDeviceStatus(c *gin.Context) {
	// 定义请求结构体（从 JSON 请求体中解析）
	var req struct {
		ID     uint `json:"id" binding:"required"`
		Status int  `json:"status"`
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	db := database.GetDB()

	// 查找对应的设备
	var device models.Device
	if err := db.First(&device, req.ID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	// 更新状态
	device.Status = req.Status
	if err := db.Save(&device).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update device status"})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"message": "Device status updated successfully",
		"device":  device,
	})
}

// CreateDevice godoc
// @Summary 创建新设备
// @Description 创建一个新的设备记录
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "创建设备请求" SchemaExample({"name": "设备名称", "status": 1, "position": "设备位置"})
// @Success 200 {object} map[string]interface{} "创建成功响应"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /device [post]
func CreateDevice(c *gin.Context) {
	// 定义请求结构体
	var req struct {
		Name     string `json:"name" binding:"required"`     // 设备名称，不能为空但可以为零值
		Status   int    `json:"status"`                      // 设备状态
		Position string `json:"position" binding:"required"` // 设备位置
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	db := database.GetDB()

	// 创建设备对象
	device := models.Device{
		Name:     req.Name,
		Status:   req.Status,
		Position: req.Position,
	}

	// 插入数据库
	if err := db.Create(&device).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create device"})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"message": "Device created successfully",
		"device": gin.H{
			"id":       device.ID,
			"name":     device.Name,
			"status":   device.Status,
			"position": device.Position,
		},
	})
}

// PaginationDevices godoc
// @Summary 设备分页查询
// @Description 根据条件对设备进行分页查询
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param name query string false "设备名称(模糊查询)"
// @Param position query string false "设备位置(模糊查询)"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "分页查询结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Router /devices/pagination [get]
func GetPagination(c *gin.Context) {
	var devices []models.Device

	// 获取查询参数
	name := c.DefaultQuery("name", "")
	position := c.DefaultQuery("position", "")

	// 构建基础查询
	query := database.DB.Model(&models.Device{})

	// 动态添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if position != "" {
		query = query.Where("position LIKE ?", "%"+position+"%")
	}

	// 执行查询获取所有匹配数据
	if err := query.Find(&devices).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid query conditions"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	p, err := strconv.Atoi(page)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid page number"})
		return
	}
	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid pagesize number"})
		return
	}
	result := PaginationDevices(devices, p, ps)

	c.JSON(200, gin.H{
		"data":  result,
		"total": len(devices),
	})
}

func PaginationDevices(data []models.Device, page, pageSize int) []models.Device {
	start := (page - 1) * pageSize
	if start >= len(data) {
		return []models.Device{}
	}
	end := start + pageSize
	if end > len(data) {
		end = len(data)
	}
	return data[start:end]
}