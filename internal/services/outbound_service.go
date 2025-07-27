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
		OrderNo:         orderNo,
		DeliveryAddress: req.DeliveryAddress,
		CarNumber:       req.CarNumber,
		DriverName:      req.DriverName,
		DriverPhone:     req.DriverPhone,
		TotalAmount:     totalAmount,
		Status:          "completed",
		Notes:           req.Notes,
		CreatedBy:       createdBy,
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
	updates := map[string]interface{}{"status": status}
	return s.outboundRepo.UpdateFields(id, updates)
}

// UpdateCustomerName 显式更新客户名称
func (s *OutboundService) UpdateCustomerName(id uint, customerName string) error {
	updates := map[string]interface{}{"delivery_address": customerName}
	return s.outboundRepo.UpdateFields(id, updates)
}

// UpdateNotes 显式更新备注
func (s *OutboundService) UpdateNotes(id uint, notes string) error {
	updates := map[string]interface{}{"notes": notes}
	return s.outboundRepo.UpdateFields(id, updates)
}

// UpdateOrder 显式更新订单字段
func (s *OutboundService) UpdateOrder(id uint, updates map[string]interface{}) error {
	return s.outboundRepo.UpdateFields(id, updates)
}

func (s *OutboundService) Delete(id uint) error {
	return s.outboundRepo.Delete(id)
}

// UpdateOrderComplete 完整更新出库订单（包括订单项）
func (s *OutboundService) UpdateOrderComplete(id uint, req *models.UpdateOutboundOrderRequest) error {
	// 检查订单是否存在
	order, err := s.outboundRepo.GetByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	// 如果需要更新订单项，先处理库存恢复
	if len(req.Items) > 0 {
		// 获取当前订单项以便恢复库存
		currentItems, err := s.outboundRepo.GetRawItemsByOrderID(id)
		if err != nil {
			return err
		}

		// 恢复当前订单项的库存
		for _, item := range currentItems {
			if err := s.inventoryRepo.UpdateWeight(item.CategoryID, item.Weight, true); err != nil {
				return err
			}
		}

		// 删除当前所有订单项
		if err := s.outboundRepo.DeleteItemsByOrderID(id); err != nil {
			return err
		}

		// 处理新的订单项并计算新的总金额
		var totalAmount float64
		for _, reqItem := range req.Items {
			// 检查库存是否足够
			inventory, err := s.inventoryRepo.GetByCategoryID(reqItem.CategoryID)
			if err != nil {
				return errors.New("inventory not found for category")
			}

			if inventory.CurrentWeightKg < reqItem.Weight {
				return errors.New("insufficient inventory for category")
			}

			// 创建新订单项
			subTotal := reqItem.Weight * reqItem.UnitPrice
			totalAmount += subTotal

			newItem := &models.OutboundOrderItem{
				OrderID:    id,
				CategoryID: reqItem.CategoryID,
				Weight:     reqItem.Weight,
				UnitPrice:  reqItem.UnitPrice,
				SubTotal:   subTotal,
			}

			if err := s.outboundRepo.CreateItem(newItem); err != nil {
				return err
			}

			// 更新库存（减少）
			if err := s.inventoryRepo.UpdateWeight(reqItem.CategoryID, reqItem.Weight, false); err != nil {
				return err
			}
		}

		// 更新订单总金额
		order.TotalAmount = totalAmount
	}

	// 更新订单基本信息
	updates := make(map[string]interface{})
	if req.DeliveryAddress != "" {
		updates["delivery_address"] = req.DeliveryAddress
		order.DeliveryAddress = req.DeliveryAddress
	}
	if req.CarNumber != "" {
		updates["car_number"] = req.CarNumber
		order.CarNumber = req.CarNumber
	}
	if req.DriverName != "" {
		updates["driver_name"] = req.DriverName
		order.DriverName = req.DriverName
	}
	if req.DriverPhone != "" {
		updates["driver_phone"] = req.DriverPhone
		order.DriverPhone = req.DriverPhone
	}
	if req.Status != "" {
		updates["status"] = req.Status
		order.Status = req.Status
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
		order.Notes = req.Notes
	}
	if len(req.Items) > 0 {
		updates["total_amount"] = order.TotalAmount
	}

	// 执行更新
	if len(updates) > 0 {
		if err := s.outboundRepo.UpdateFields(id, updates); err != nil {
			return err
		}
	}

	return nil
}

// UpdateOrderBasic 仅更新订单基本信息（不包括订单项）
func (s *OutboundService) UpdateOrderBasic(id uint, req *models.UpdateOutboundOrderRequest) error {
	updates := make(map[string]interface{})

	if req.DeliveryAddress != "" {
		updates["delivery_address"] = req.DeliveryAddress
	}
	if req.CarNumber != "" {
		updates["car_number"] = req.CarNumber
	}
	if req.DriverName != "" {
		updates["driver_name"] = req.DriverName
	}
	if req.DriverPhone != "" {
		updates["driver_phone"] = req.DriverPhone
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	return s.outboundRepo.UpdateFields(id, updates)
}
