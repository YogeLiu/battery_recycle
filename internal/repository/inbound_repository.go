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

// InboundRepository 入库订单数据仓库 (不再是接口实现)
type InboundRepository struct {
	db *gorm.DB
}

// NewInboundRepository 创建入库仓库实例
func NewInboundRepository(db *gorm.DB) *InboundRepository {
	return &InboundRepository{db: db}
}

// Create 创建入库订单
func (r *InboundRepository) Create(order *models.InboundOrder) error {
	return r.db.Create(order).Error
}

// CreateItem 创建入库订单项
func (r *InboundRepository) CreateItem(item *models.InboundOrderItem) error {
	return r.db.Create(item).Error
}

// GetByID 根据ID获取入库订单
func (r *InboundRepository) GetByID(id uint) (*models.InboundOrder, error) {
	var order models.InboundOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAll 获取所有入库订单 (分页)
func (r *InboundRepository) GetAll(limit, offset int) ([]models.InboundOrder, int64, error) {
	var orders []models.InboundOrder
	var total int64

	// Get total count
	r.db.Model(&models.InboundOrder{}).Count(&total)

	// Get orders with pagination
	err := r.db.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

// GetAllWithConditions 根据条件获取入库订单 (支持筛选和分页)
func (r *InboundRepository) GetAllWithConditions(req *models.GetInboundOrderRequest) ([]models.InboundOrder, int64, error) {
	query := r.db.Model(&models.InboundOrder{})

	// 应用筛选条件
	if req.Supplier != "" {
		query = query.Where("supplier_name LIKE ?", "%"+req.Supplier+"%")
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at >= ? AND created_at <= ?", req.StartDate, req.EndDate)
	}
	if req.StartDate != "" && req.EndDate == "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}

	query = query.Where("is_deleted = 0")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页和排序
	var orders []models.InboundOrder
	offset := (req.Page - 1) * req.PageSize
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	err := query.Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&orders).Error

	return orders, total, err
}

// UpdateStatus 显式更新订单状态
func (r *InboundRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateSupplierName 显式更新供应商名称
func (r *InboundRepository) UpdateSupplierName(id uint, supplierName string) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Update("supplier_name", supplierName).Error
}

// UpdateNotes 显式更新备注
func (r *InboundRepository) UpdateNotes(id uint, notes string) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Update("notes", notes).Error
}

// UpdateTotalAmount 显式更新总金额
func (r *InboundRepository) UpdateTotalAmount(id uint, totalAmount float64) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Update("total_amount", totalAmount).Error
}

// UpdateFields 显式更新指定字段
func (r *InboundRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除入库订单
func (r *InboundRepository) Delete(id uint) error {
	return r.db.Model(&models.InboundOrder{}).Where("id = ?", id).Update("is_deleted", 1).Error
}

// GetItemsByOrderID 根据订单ID获取订单详细条目 (包含分类名称)
func (r *InboundRepository) GetItemsByOrderID(orderID uint) ([]models.InboundOrderDetailDTO, error) {
	var result []models.InboundOrderDetailDTO

	err := r.db.Table("inbound_order_items as i").
		Select(`
			i.category_id,
			COALESCE(c.name, '未知分类') as category_name,
			i.gross_weight,
			i.tare_weight,
			i.net_weight,
			i.unit_price,
			i.sub_total
		`).
		Joins("LEFT JOIN battery_categories c ON i.category_id = c.id").
		Where("i.order_id = ?", orderID).
		Scan(&result).Error

	return result, err
}

// GenerateOrderNo 生成订单号 (并发安全：纳秒时间戳+随机数)
func (r *InboundRepository) GenerateOrderNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	// 使用纳秒时间戳确保时间唯一性
	nanoStr := fmt.Sprintf("%05d", now.Nanosecond())

	// 生成4位随机数字进一步保证唯一性
	randomNum := rand.Intn(10000) // 0-9999
	randomStr := fmt.Sprintf("%04d", randomNum)

	// 格式：IN-20240101-123456789-1234
	orderNo := fmt.Sprintf("IN-%s-%s-%s", dateStr, nanoStr, randomStr)
	return orderNo, nil
}
