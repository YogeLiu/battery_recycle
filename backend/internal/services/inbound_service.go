package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

type inboundService struct {
	inboundRepo   repository.InboundRepository
	inventoryRepo repository.InventoryRepository
}

func NewInboundService(inboundRepo repository.InboundRepository, inventoryRepo repository.InventoryRepository) InboundService {
	return &inboundService{
		inboundRepo:   inboundRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *inboundService) Create(req *models.CreateInboundOrderRequest, createdBy uint) (*models.InboundOrder, error) {
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

func (s *inboundService) GetByID(id uint) (*models.InboundOrder, error) {
	return s.inboundRepo.GetByID(id)
}

func (s *inboundService) GetAll(limit, offset int) ([]models.InboundOrder, int64, error) {
	return s.inboundRepo.GetAll(limit, offset)
}

func (s *inboundService) Update(order *models.InboundOrder) error {
	return s.inboundRepo.Update(order)
}

func (s *inboundService) Delete(id uint) error {
	return s.inboundRepo.Delete(id)
}
