package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
	"errors"
)

// OutboundService 出库服务 (不再使用接口)
type OutboundService struct {
	outboundRepo  *repository.OutboundRepository
	inventoryRepo *repository.InventoryRepository
}

// NewOutboundService 创建出库服务实例
func NewOutboundService(outboundRepo *repository.OutboundRepository, inventoryRepo *repository.InventoryRepository) *OutboundService {
	return &OutboundService{
		outboundRepo:  outboundRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *OutboundService) Create(req *models.CreateOutboundOrderRequest, createdBy uint) (*models.OutboundOrder, error) {
	// Generate order number
	orderNo, err := s.outboundRepo.GenerateOrderNo()
	if err != nil {
		return nil, err
	}

	// Check inventory and calculate totals
	var totalAmount float64
	var orderItems []models.OutboundOrderItem

	for _, reqItem := range req.Items {
		// Check inventory availability
		inventory, err := s.inventoryRepo.GetByCategoryID(reqItem.CategoryID)
		if err != nil {
			return nil, errors.New("inventory not found for category")
		}

		if inventory.CurrentWeightKg < reqItem.Weight {
			return nil, errors.New("insufficient inventory")
		}

		subTotal := reqItem.Weight * reqItem.UnitPrice
		totalAmount += subTotal

		orderItem := models.OutboundOrderItem{
			CategoryID: reqItem.CategoryID,
			Weight:     reqItem.Weight,
			UnitPrice:  reqItem.UnitPrice,
			SubTotal:   subTotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &models.OutboundOrder{
		OrderNo:      orderNo,
		CustomerName: req.CustomerName,
		TotalAmount:  totalAmount,
		Status:       "completed",
		Notes:        req.Notes,
		CreatedBy:    createdBy,
	}

	if err := s.outboundRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items and update inventory
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.outboundRepo.CreateItem(&orderItems[i]); err != nil {
			return nil, err
		}

		// Update inventory
		if err := s.inventoryRepo.UpdateWeight(orderItems[i].CategoryID, orderItems[i].Weight, false); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// GetByID 根据ID获取出库订单详情 (包含详细条目)
func (s *OutboundService) GetByID(id uint) (*models.GetOutboundOrderDetailResp, error) {
	order, err := s.outboundRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := s.outboundRepo.GetItemsByOrderID(id)
	if err != nil {
		return nil, err
	}

	return &models.GetOutboundOrderDetailResp{Order: *order, Detail: items}, nil
}

// GetAll 获取所有出库订单 (支持条件筛选和分页)
func (s *OutboundService) GetAll(req *models.GetOutboundOrderRequest) ([]models.OutboundOrder, int64, error) {
	return s.outboundRepo.GetAllWithConditions(req)
}

// GetAllSimple 获取所有出库订单 (简单分页，保持向后兼容)
func (s *OutboundService) GetAllSimple(limit, offset int) ([]models.OutboundOrder, int64, error) {
	return s.outboundRepo.GetAll(limit, offset)
}

// UpdateStatus 显式更新订单状态
func (s *OutboundService) UpdateStatus(id uint, status string) error {
	return s.outboundRepo.UpdateStatus(id, status)
}

// UpdateCustomerName 显式更新客户名称
func (s *OutboundService) UpdateCustomerName(id uint, customerName string) error {
	return s.outboundRepo.UpdateCustomerName(id, customerName)
}

// UpdateNotes 显式更新备注
func (s *OutboundService) UpdateNotes(id uint, notes string) error {
	return s.outboundRepo.UpdateNotes(id, notes)
}

// UpdateOrder 显式更新订单字段
func (s *OutboundService) UpdateOrder(id uint, updates map[string]interface{}) error {
	return s.outboundRepo.UpdateFields(id, updates)
}

func (s *OutboundService) Delete(id uint) error {
	return s.outboundRepo.Delete(id)
}
