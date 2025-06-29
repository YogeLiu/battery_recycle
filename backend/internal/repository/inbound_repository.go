package repository

import (
	"battery-erp-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type inboundRepository struct {
	db *gorm.DB
}

func NewInboundRepository(db *gorm.DB) InboundRepository {
	return &inboundRepository{db: db}
}

func (r *inboundRepository) Create(order *models.InboundOrder) error {
	return r.db.Create(order).Error
}

func (r *inboundRepository) CreateItem(item *models.InboundOrderItem) error {
	return r.db.Create(item).Error
}

func (r *inboundRepository) GetByID(id uint) (*models.InboundOrder, error) {
	var order models.InboundOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *inboundRepository) GetAll(limit, offset int) ([]models.InboundOrder, int64, error) {
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

func (r *inboundRepository) Update(order *models.InboundOrder) error {
	return r.db.Save(order).Error
}

func (r *inboundRepository) Delete(id uint) error {
	return r.db.Delete(&models.InboundOrder{}, id).Error
}

func (r *inboundRepository) GenerateOrderNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	var count int64
	r.db.Model(&models.InboundOrder{}).
		Where("DATE(created_at) = ?", now.Format("2006-01-02")).
		Count(&count)

	orderNo := fmt.Sprintf("IN-%s-%04d", dateStr, count+1)
	return orderNo, nil
}
