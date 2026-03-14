package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	host := "mysql-0.mysql.default.svc.cluster.local"
	port := "3306"
	user := "root"
	password := "root"
	dbname := "gorder_v2"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&allowNativePasswords=true",
		user, password, host, port, dbname,
	)

	fmt.Printf("Testing GORM connection to %s:%s...\n", host, port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("GORM connection failed: %v\n", err)
		return
	}
	fmt.Printf("GORM connection success: %s\n", dsn)
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Get sql.DB failed: %v\n", err)
		return
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		return
	}

	fmt.Println("GORM connection OK!")

	// 查询数据
	var result []map[string]interface{}
	if err := db.Table("o_stock").Scan(&result).Error; err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return
	}

	fmt.Printf("Query result: %v\n", result)
}
