package repository

import (
	"battery-erp-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type outboundRepository struct {
	db *gorm.DB
}

func NewOutboundRepository(db *gorm.DB) OutboundRepository {
	return &outboundRepository{db: db}
}

func (r *outboundRepository) Create(order *models.OutboundOrder) error {
	return r.db.Create(order).Error
}

func (r *outboundRepository) CreateItem(item *models.OutboundOrderItem) error {
	return r.db.Create(item).Error
}

func (r *outboundRepository) GetByID(id uint) (*models.OutboundOrder, error) {
	var order models.OutboundOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *outboundRepository) GetAll(limit, offset int) ([]models.OutboundOrder, int64, error) {
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

func (r *outboundRepository) Update(order *models.OutboundOrder) error {
	return r.db.Save(order).Error
}

func (r *outboundRepository) Delete(id uint) error {
	return r.db.Delete(&models.OutboundOrder{}, id).Error
}

func (r *outboundRepository) GenerateOrderNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	var count int64
	r.db.Model(&models.OutboundOrder{}).
		Where("DATE(created_at) = ?", now.Format("2006-01-02")).
		Count(&count)

	orderNo := fmt.Sprintf("OUT-%s-%04d", dateStr, count+1)
	return orderNo, nil
}
