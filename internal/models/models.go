package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	RealName  string    `json:"real_name" gorm:"size:100;not null"`
	Role      string    `json:"role" gorm:"size:20;not null;default:'normal'"` // 'super_admin' or 'normal'
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (User) TableName() string {
	return "users"
}

// BatteryCategory represents a battery category/type
type BatteryCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:255"`
	UnitPrice   float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"` // Price per kg
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (BatteryCategory) TableName() string {
	return "battery_categories"
}

// InboundOrder represents a purchase/inbound order
type InboundOrder struct {
	ID           uint      `json:"id" gorm:"primaryKey"`                               // 订单ID
	OrderNo      string    `json:"order_no" gorm:"uniqueIndex;size:50;not null"`       // 订单号
	SupplierName string    `json:"supplier_name" gorm:"size:100;not null"`             // 供应商名称
	TotalAmount  float64   `json:"total_amount" gorm:"type:decimal(15,2);not null"`    // 总金额
	Status       string    `json:"status" gorm:"size:20;not null;default:'completed'"` // 'completed', 'cancelled'
	Notes        string    `json:"notes" gorm:"type:text"`                             // 备注
	CreatedBy    uint      `json:"created_by" gorm:"not null"`                         // 创建人
	IsDeleted    int       `json:"is_deleted" gorm:"default:0"`                        // 是否删除
	CreatedAt    time.Time `json:"created_at"`                                         // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`                                         // 更新时间
}

// TableName sets the insert table name for this struct type
func (InboundOrder) TableName() string {
	return "inbound_orders"
}

// InboundOrderItem represents items in an inbound order
type InboundOrderItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`                            // 订单项ID
	OrderID     uint      `json:"order_id" gorm:"not null"`                        // 订单ID
	CategoryID  uint      `json:"category_id" gorm:"not null"`                     // 电池类型ID
	GrossWeight float64   `json:"gross_weight" gorm:"type:decimal(10,3);not null"` // kg
	TareWeight  float64   `json:"tare_weight" gorm:"type:decimal(10,3);not null"`  // kg
	NetWeight   float64   `json:"net_weight" gorm:"type:decimal(10,3);not null"`   // kg
	UnitPrice   float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"`   // 单价
	SubTotal    float64   `json:"sub_total" gorm:"type:decimal(15,2);not null"`    // 小计
	CreatedAt   time.Time `json:"created_at"`                                      // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                                      // 更新时间
}

// TableName sets the insert table name for this struct type
func (InboundOrderItem) TableName() string {
	return "inbound_order_items"
}

// OutboundOrder represents a sales/outbound order
type OutboundOrder struct {
	ID              uint      `json:"id" gorm:"primaryKey"`                               // 主键
	OrderNo         string    `json:"order_no" gorm:"size:50;uniqueIndex;not null"`       // 订单号 YYYYMMDD999999
	DeliveryAddress string    `json:"delivery_address" gorm:"size:255;not null"`          // 送货地
	CarNumber       string    `json:"car_number" gorm:"size:50;not null"`                 // 车号
	DriverName      string    `json:"driver_name" gorm:"size:50;not null"`                // 司机姓名
	DriverPhone     string    `json:"driver_phone" gorm:"size:20;not null"`               // 司机手机号
	TotalAmount     float64   `json:"total_amount" gorm:"type:decimal(15,2);not null"`    // 总金额
	Status          string    `json:"status" gorm:"size:20;not null;default:'completed'"` // 'completed', 'cancelled'
	Notes           string    `json:"notes" gorm:"type:text"`                             // 备注
	CreatedBy       uint      `json:"created_by" gorm:"not null"`                         // 创建人
	IsDeleted       int       `json:"is_deleted" gorm:"default:0"`                        // 是否删除
	CreatedAt       time.Time `json:"created_at"`                                         // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`                                         // 更新时间
}

// TableName sets the insert table name for this struct type
func (OutboundOrder) TableName() string {
	return "outbound_orders"
}

// OutboundOrderItem represents items in an outbound order
type OutboundOrderItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	OrderID    uint      `json:"order_id" gorm:"not null"`
	CategoryID uint      `json:"category_id" gorm:"not null"`
	Weight     float64   `json:"weight" gorm:"type:decimal(10,3);not null"`     // kg
	UnitPrice  float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"` // Price per kg
	SubTotal   float64   `json:"sub_total" gorm:"type:decimal(15,2);not null"`  // Weight * unit price
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (OutboundOrderItem) TableName() string {
	return "outbound_order_items"
}

// Inventory represents current inventory for each battery category
type Inventory struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	CategoryID      uint       `json:"category_id" gorm:"uniqueIndex;not null"`
	CurrentWeightKg float64    `json:"current_weight_kg" gorm:"type:decimal(12,3);not null;default:0"`
	LastInboundAt   *time.Time `json:"last_inbound_at"`
	LastOutboundAt  *time.Time `json:"last_outbound_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Inventory) TableName() string {
	return "inventories"
}

// Common response structures
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ReportSummary represents report summary data
type ReportSummary struct {
	TotalInventoryWeight float64           `json:"total_inventory_weight"`
	InventoryCount       int               `json:"inventory_count"`
	InventoryDetails     []InventoryDetail `json:"inventory_details"`
	InboundStats         *OrderStats       `json:"inbound_stats,omitempty"`
	OutboundStats        *OrderStats       `json:"outbound_stats,omitempty"`
	ReportGeneratedAt    time.Time         `json:"report_generated_at"`
	DateRange            *DateRange        `json:"date_range,omitempty"`
}

// InventoryDetail represents inventory detail in report
type InventoryDetail struct {
	CategoryID     uint       `json:"category_id"`
	CategoryName   string     `json:"category_name"`
	CurrentWeight  float64    `json:"current_weight"`
	LastInboundAt  *time.Time `json:"last_inbound_at"`
	LastOutboundAt *time.Time `json:"last_outbound_at"`
}

// OrderStats represents order statistics
type OrderStats struct {
	TotalOrders int64   `json:"total_orders"`
	TotalAmount float64 `json:"total_amount"`
	TotalWeight float64 `json:"total_weight"`
	AvgAmount   float64 `json:"avg_amount"`
}

// DateRange represents date range for reports
type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// CreateInboundOrderRequest represents request to create inbound order
type CreateInboundOrderRequest struct {
	SupplierName string                   `json:"supplier_name" binding:"required"`
	Notes        string                   `json:"notes"`
	Items        []CreateInboundOrderItem `json:"items" binding:"required,dive"`
}

// CreateInboundOrderItem represents item in create inbound order request
type CreateInboundOrderItem struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	GrossWeight float64 `json:"gross_weight" binding:"required,gt=0"`
	TareWeight  float64 `json:"tare_weight" binding:"required,gte=0"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
}

// CreateOutboundOrderRequest represents request to create outbound order
type CreateOutboundOrderRequest struct {
	DeliveryAddress string                    `json:"delivery_address" binding:"required"`
	CarNumber       string                    `json:"car_number" binding:"required"`
	DriverName      string                    `json:"driver_name" binding:"required"`
	DriverPhone     string                    `json:"driver_phone" binding:"required"`
	Notes           string                    `json:"notes"`
	Items           []CreateOutboundOrderItem `json:"items" binding:"required,dive"`
}

// CreateOutboundOrderItem represents item in create outbound order request
type CreateOutboundOrderItem struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Weight     float64 `json:"weight" binding:"required,gt=0"`
	UnitPrice  float64 `json:"unit_price" binding:"required,gt=0"`
}

// Business error codes
const (
	// Success codes
	CodeSuccess = 0

	// Client error codes (4xx equivalent)
	CodeBadRequest       = 40000 // 400 - Bad Request
	CodeUnauthorized     = 40100 // 401 - Unauthorized
	CodeForbidden        = 40300 // 403 - Forbidden
	CodeNotFound         = 40400 // 404 - Not Found
	CodeMethodNotAllowed = 40500 // 405 - Method Not Allowed
	CodeConflict         = 40900 // 409 - Conflict

	// Server error codes (5xx equivalent)
	CodeInternalError      = 50000 // 500 - Internal Server Error
	CodeBadGateway         = 50200 // 502 - Bad Gateway
	CodeServiceUnavailable = 50300 // 503 - Service Unavailable
)

type UpdateCategoryRequest struct {
	ID          uint    `json:"-"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
}
type GetInboundOrderRequest struct {
	Page      int    `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize  int    `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"`
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
	Supplier  string `json:"supplier" form:"supplier"`
}

type GetInboundOrderResponse struct {
	Orders []InboundOrder `json:"orders"`
	Total  int64          `json:"total"`
}

type InboundOrderDetailDTO struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	GrossWeight  float64 `json:"gross_weight"`
	TareWeight   float64 `json:"tare_weight"`
	NetWeight    float64 `json:"net_weight"`
	UnitPrice    float64 `json:"unit_price"`
	SubTotal     float64 `json:"sub_total"`
}

type GetInboudOrderDetailResp struct {
	Order  InboundOrder            `json:"order"`
	Detail []InboundOrderDetailDTO `json:"detail"`
}

// 出库订单相关模型
type OutboundOrderDetailDTO struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Weight       float64 `json:"weight"`
	UnitPrice    float64 `json:"unit_price"`
	SubTotal     float64 `json:"sub_total"`
}

type GetOutboundOrderDetailResp struct {
	Order  OutboundOrder            `json:"order"`
	Detail []OutboundOrderDetailDTO `json:"detail"`
}

type GetOutboundOrderRequest struct {
	Page      int    `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize  int    `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100"`
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
	Customer  string `json:"customer" form:"customer"`
}

type GetOutboundOrderResponse struct {
	Orders []OutboundOrder `json:"orders"`
	Total  int64           `json:"total"`
}

// UpdateOutboundOrderRequest represents request to update outbound order
type UpdateOutboundOrderRequest struct {
	DeliveryAddress string                    `json:"delivery_address"`
	CarNumber       string                    `json:"car_number"`
	DriverName      string                    `json:"driver_name"`
	DriverPhone     string                    `json:"driver_phone"`
	Status          string                    `json:"status"`
	Notes           string                    `json:"notes"`
	Items           []UpdateOutboundOrderItem `json:"items,omitempty"`
}

// UpdateOutboundOrderItem represents item in update outbound order request
type UpdateOutboundOrderItem struct {
	ID         uint    `json:"id,omitempty"` // 如果有ID则是更新，没有则是新增
	CategoryID uint    `json:"category_id" binding:"required"`
	Weight     float64 `json:"weight" binding:"required,gt=0"`
	UnitPrice  float64 `json:"unit_price" binding:"required,gt=0"`
	Action     string  `json:"action,omitempty"` // "add", "update", "delete"
}
