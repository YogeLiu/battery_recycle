package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
	"errors"
)

type outboundService struct {
	outboundRepo  repository.OutboundRepository
	inventoryRepo repository.InventoryRepository
}

func NewOutboundService(outboundRepo repository.OutboundRepository, inventoryRepo repository.InventoryRepository) OutboundService {
	return &outboundService{
		outboundRepo:  outboundRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *outboundService) Create(req *models.CreateOutboundOrderRequest, createdBy uint) (*models.OutboundOrder, error) {
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

func (s *outboundService) GetByID(id uint) (*models.OutboundOrder, error) {
	return s.outboundRepo.GetByID(id)
}

func (s *outboundService) GetAll(limit, offset int) ([]models.OutboundOrder, int64, error) {
	return s.outboundRepo.GetAll(limit, offset)
}

func (s *outboundService) Update(order *models.OutboundOrder) error {
	return s.outboundRepo.Update(order)
}

func (s *outboundService) Delete(id uint) error {
	return s.outboundRepo.Delete(id)
}
