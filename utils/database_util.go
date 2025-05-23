package utils

import (
	"cms/config"
	"cms/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.DBUser,
		dbConfig.DBPassword,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	)

	datetimePrecision := 2
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
		return nil, err
	}

	db.AutoMigrate(&models.Category{}, &models.User{}, &models.Image{}, &models.Tag{}, &models.Article{}, &models.Dict{})

	return db, nil
}
