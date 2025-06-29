package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

// InboundService 入库服务 (不再使用接口)
type InboundService struct {
	inboundRepo   *repository.InboundRepository
	inventoryRepo *repository.InventoryRepository
}

// NewInboundService 创建入库服务实例
func NewInboundService(inboundRepo *repository.InboundRepository, inventoryRepo *repository.InventoryRepository) *InboundService {
	return &InboundService{
		inboundRepo:   inboundRepo,
		inventoryRepo: inventoryRepo,
	}
}

// Create 创建入库订单
func (s *InboundService) Create(req *models.CreateInboundOrderRequest, createdBy uint) (*models.InboundOrder, error) {
	// Generate order number
	orderNo, err := s.inboundRepo.GenerateOrderNo()
	if err != nil {
		return nil, err
	}

	// Calculate item totals
	var totalAmount float64
	var orderItems []models.InboundOrderItem

	for _, reqItem := range req.Items {
		netWeight := reqItem.GrossWeight - reqItem.TareWeight
		subTotal := netWeight * reqItem.UnitPrice
		totalAmount += subTotal

		orderItem := models.InboundOrderItem{
			CategoryID:  reqItem.CategoryID,
			GrossWeight: reqItem.GrossWeight,
			TareWeight:  reqItem.TareWeight,
			NetWeight:   netWeight,
			UnitPrice:   reqItem.UnitPrice,
			SubTotal:    subTotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &models.InboundOrder{
		OrderNo:      orderNo,
		SupplierName: req.SupplierName,
		TotalAmount:  totalAmount,
		Status:       "completed",
		Notes:        req.Notes,
		CreatedBy:    createdBy,
	}

	if err := s.inboundRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.inboundRepo.CreateItem(&orderItems[i]); err != nil {
			return nil, err
		}

		// Update inventory
		if err := s.inventoryRepo.UpdateWeight(orderItems[i].CategoryID, orderItems[i].NetWeight, true); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// GetByID 根据ID获取入库订单
func (s *InboundService) GetByID(id uint) (*models.GetInboudOrderDetailResp, error) {
	order, err := s.inboundRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := s.inboundRepo.GetItemsByOrderID(id)
	if err != nil {
		return nil, err
	}

	return &models.GetInboudOrderDetailResp{Order: *order, Detail: items}, nil
}

// GetAll 获取所有入库订单 (支持条件筛选和分页)
func (s *InboundService) GetAll(req *models.GetInboundOrderRequest) ([]models.InboundOrder, int64, error) {
	return s.inboundRepo.GetAllWithConditions(req)
}

// UpdateStatus 显式更新订单状态
func (s *InboundService) UpdateStatus(id uint, status string) error {
	return s.inboundRepo.UpdateStatus(id, status)
}

// UpdateSupplierName 显式更新供应商名称
func (s *InboundService) UpdateSupplierName(id uint, supplierName string) error {
	return s.inboundRepo.UpdateSupplierName(id, supplierName)
}

// UpdateNotes 显式更新备注
func (s *InboundService) UpdateNotes(id uint, notes string) error {
	return s.inboundRepo.UpdateNotes(id, notes)
}

// UpdateOrder 显式更新订单字段
func (s *InboundService) UpdateOrder(id uint, updates map[string]interface{}) error {
	return s.inboundRepo.UpdateFields(id, updates)
}

// Delete 删除入库订单
func (s *InboundService) Delete(id uint) error {
	return s.inboundRepo.Delete(id)
}
