package repository

import (
	"battery-erp-backend/internal/models"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// OutboundRepository 出库订单数据仓库 (不再是接口实现)
type OutboundRepository struct {
	db *gorm.DB
}

// NewOutboundRepository 创建出库仓库实例
func NewOutboundRepository(db *gorm.DB) *OutboundRepository {
	return &OutboundRepository{db: db}
}

// Create 创建出库订单
func (r *OutboundRepository) Create(order *models.OutboundOrder) error {
	return r.db.Create(order).Error
}

// CreateItem 创建出库订单项
func (r *OutboundRepository) CreateItem(item *models.OutboundOrderItem) error {
	return r.db.Create(item).Error
}

// GetByID 根据ID获取出库订单
func (r *OutboundRepository) GetByID(id uint) (*models.OutboundOrder, error) {
	var order models.OutboundOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAll 获取所有出库订单 (分页)
func (r *OutboundRepository) GetAll(limit, offset int) ([]models.OutboundOrder, int64, error) {
	var orders []models.OutboundOrder
	var total int64

	// Get total count
	r.db.Model(&models.OutboundOrder{}).Count(&total)

	// Get orders with pagination
	err := r.db.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

// UpdateStatus 显式更新订单状态
func (r *OutboundRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.OutboundOrder{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateCustomerName 显式更新客户名称
func (r *OutboundRepository) UpdateCustomerName(id uint, customerName string) error {
	return r.db.Model(&models.OutboundOrder{}).Where("id = ?", id).Update("customer_name", customerName).Error
}

// UpdateNotes 显式更新备注
func (r *OutboundRepository) UpdateNotes(id uint, notes string) error {
	return r.db.Model(&models.OutboundOrder{}).Where("id = ?", id).Update("notes", notes).Error
}

// UpdateTotalAmount 显式更新总金额
func (r *OutboundRepository) UpdateTotalAmount(id uint, totalAmount float64) error {
	return r.db.Model(&models.OutboundOrder{}).Where("id = ?", id).Update("total_amount", totalAmount).Error
}

// UpdateFields 显式更新指定字段
func (r *OutboundRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.OutboundOrder{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除出库订单
func (r *OutboundRepository) Delete(id uint) error {
	return r.db.Delete(&models.OutboundOrder{}, id).Error
}

// GetItemsByOrderID 根据订单ID获取出库订单详细条目 (包含分类名称)
func (r *OutboundRepository) GetItemsByOrderID(orderID uint) ([]models.OutboundOrderDetailDTO, error) {
	var result []models.OutboundOrderDetailDTO

	err := r.db.Table("outbound_order_items as o").
		Select(`
			o.category_id,
			COALESCE(c.name, '未知分类') as category_name,
			o.weight,
			o.unit_price,
			o.sub_total
		`).
		Joins("LEFT JOIN battery_categories c ON o.category_id = c.id").
		Where("o.order_id = ?", orderID).
		Scan(&result).Error

	return result, err
}

// GetAllWithConditions 根据条件获取出库订单 (支持筛选和分页)
func (r *OutboundRepository) GetAllWithConditions(req *models.GetOutboundOrderRequest) ([]models.OutboundOrder, int64, error) {
	query := r.db.Model(&models.OutboundOrder{})

	// 应用筛选条件
	if req.Customer != "" {
		query = query.Where("customer_name LIKE ?", "%"+req.Customer+"%")
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at >= ? AND created_at <= ?", req.StartDate, req.EndDate)
	}
	if req.StartDate != "" && req.EndDate == "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页和排序
	var orders []models.OutboundOrder
	offset := (req.Page - 1) * req.PageSize
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&orders).Error

	return orders, total, err
}

// GenerateOrderNo 生成订单号 (并发安全：纳秒时间戳+随机数)
func (r *OutboundRepository) GenerateOrderNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	// 使用纳秒时间戳确保时间唯一性
	nanoStr := fmt.Sprintf("%05d", now.Nanosecond())

	// 生成4位随机数字进一步保证唯一性
	randomNum := rand.Intn(10000) // 0-9999
	randomStr := fmt.Sprintf("%04d", randomNum)

	// 格式：OUT-20240101-123456789-1234
	orderNo := fmt.Sprintf("OUT-%s-%s-%s", dateStr, nanoStr, randomStr)
	return orderNo, nil
}
