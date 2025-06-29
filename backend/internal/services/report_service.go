package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
	"time"
)

type reportService struct {
	repos *repository.Repositories
}

func NewReportService(repos *repository.Repositories) ReportService {
	return &reportService{
		repos: repos,
	}
}

func (s *reportService) GetSummary(startDate, endDate string) (*models.ReportSummary, error) {
	// 获取库存总览
	inventories, err := s.repos.Inventory.GetAll()
	if err != nil {
		return nil, err
	}

	var totalWeight float64
	var inventoryDetails []models.InventoryDetail

	for _, inv := range inventories {
		totalWeight += inv.CurrentWeightKg

		// 手动获取分类信息
		category, err := s.repos.Category.GetByID(inv.CategoryID)
		categoryName := "Unknown"
		if err == nil && category != nil {
			categoryName = category.Name
		}

		detail := models.InventoryDetail{
			CategoryID:     inv.CategoryID,
			CategoryName:   categoryName,
			CurrentWeight:  inv.CurrentWeightKg,
			LastInboundAt:  inv.LastInboundAt,
			LastOutboundAt: inv.LastOutboundAt,
		}
		inventoryDetails = append(inventoryDetails, detail)
	}

	// 构建报告摘要
	summary := &models.ReportSummary{
		TotalInventoryWeight: totalWeight,
		InventoryCount:       len(inventories),
		InventoryDetails:     inventoryDetails,
		ReportGeneratedAt:    time.Now(),
	}

	// 如果提供了日期范围，添加到报告中
	if startDate != "" && endDate != "" {
		summary.DateRange = &models.DateRange{
			StartDate: startDate,
			EndDate:   endDate,
		}

		// TODO: 可以在这里添加按日期范围的统计逻辑
		// 比如指定日期范围内的入库/出库统计
	}

	return summary, nil
}
