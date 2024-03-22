package initialize

import (
	"github.com/wilsonce/connectly-test/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
	"sync"
)

var DB *gorm.DB
var dbOnce sync.Once

func InitDB() *gorm.DB {
	dbOnce.Do(func() {
		db, err := gorm.Open(sqlite.Open("connectly_bot.db"), &gorm.Config{})
		if err != nil {
			Logger.Error(err.Error())
			os.Exit(0)
		}
		DB = db
	})
	return DB
}

func InitDao() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./dao", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance
	db := InitDB()

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.WBot{}, model.WMessage{})

	// Generate default DAO interface for those generated structs from database
	//companyGenerator := g.GenerateModelAs("company", "MyCompany"),
	//	g.ApplyBasic(
	//		g.GenerateModel("users"),
	//		companyGenerator,
	//		g.GenerateModelAs("people", "Person",
	//			gen.FieldIgnore("deleted_at"),
	//			gen.FieldNewTag("age", `json:"-"`),
	//		),
	//	)

	// Execute the generator
	g.Execute()
}

func InitModel() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	db := InitDB()
	g.UseDB(db) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.GenerateModel("vod")
	g.GenerateAllTable()

	//g.ApplyBasic(
	//	// Generate structs from all tables of current database
	//	g.GenerateAllTable()...,
	//)
	// Generate the code
	g.Execute()
}
