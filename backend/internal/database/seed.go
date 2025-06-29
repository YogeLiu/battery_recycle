package database

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"log"

	"gorm.io/gorm"
)

// SeedDatabase initializes the database with default data
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Create default admin user
	adminPassword := "admin123"
	hashedPassword, err := services.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	adminUser := &models.User{
		Username: "admin",
		Password: hashedPassword,
		RealName: "系统管理员",
		Role:     "super_admin",
		IsActive: true,
	}

	// Create or update admin user
	var existingUser models.User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(adminUser).Error; err != nil {
				return err
			}
			log.Println("Created admin user")
		} else {
			return err
		}
	} else {
		log.Println("Admin user already exists")
	}

	// Create default operator user
	operatorUser := &models.User{
		Username: "operator",
		Password: hashedPassword,
		RealName: "操作员",
		Role:     "normal",
		IsActive: true,
	}

	if err := db.Where("username = ?", "operator").First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(operatorUser).Error; err != nil {
				return err
			}
			log.Println("Created operator user")
		} else {
			return err
		}
	} else {
		log.Println("Operator user already exists")
	}

	// Create default battery categories
	categories := []models.BatteryCategory{
		{Name: "铅酸电池", Description: "废旧铅酸蓄电池", UnitPrice: 8.50, IsActive: true},
		{Name: "锂电池", Description: "废旧锂离子电池", UnitPrice: 15.00, IsActive: true},
		{Name: "镍氢电池", Description: "废旧镍氢电池", UnitPrice: 12.00, IsActive: true},
		{Name: "镍镉电池", Description: "废旧镍镉电池", UnitPrice: 10.00, IsActive: true},
		{Name: "碱性电池", Description: "废旧碱性电池", UnitPrice: 6.00, IsActive: true},
	}

	for _, category := range categories {
		var existingCategory models.BatteryCategory
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return err
				}
				log.Printf("Created category: %s", category.Name)

				// Create corresponding inventory record
				inventory := &models.Inventory{
					CategoryID:      category.ID,
					CurrentWeightKg: 0,
				}
				if err := db.Create(inventory).Error; err != nil {
					return err
				}
				log.Printf("Created inventory record for category: %s", category.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Category already exists: %s", category.Name)
		}
	}

	log.Println("Database seeding completed successfully!")
	return nil
}