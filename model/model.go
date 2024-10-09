package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sduonline-training-backend/pkg/conf"
)

var DB *gorm.DB

type AbstractModel struct {
	Tx *gorm.DB
}

func Setup() {
	dbInternal, err := gorm.Open(mysql.Open(conf.Conf.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	DB = dbInternal
}

func Database_initialization() error {
	sqlDB, err := DB.DB()
	rows, err := sqlDB.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	// 删除所有表
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			fmt.Println("Error scanning table name:", err)
			continue
		}

		result := DB.Exec("DROP TABLE IF EXISTS " + tableName)
		if result.Error != nil {
			return result.Error
		}
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&User{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&Task{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&Homework{}); err != nil {
		return err
	}
	return nil
}
