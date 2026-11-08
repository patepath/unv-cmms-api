package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"size:100;not null" json:"name"`
	Group string `gorm:"size:50;column:group" json:"group"`
}

type Brand struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

type ServiceRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AssetID     uint           `gorm:"not null" json:"assetId"`
	ServiceDate time.Time      `json:"serviceDate,omitempty"`
	Details     string         `gorm:"type:text" json:"details"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Building struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"size:20;column:code" json:"code"`
	Name      string    `gorm:"size:255;not null;column:name" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

type Asset struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`

	CategoryID *uint    `gorm:"NULL" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID;" json:"category"`

	Location string `gorm:"size:255;column:location" json:"location"`
	Floor    string `gorm:"size:20;column:floor" json:"floor"`
	AssetNo  string `gorm:"size:100;column:asset_no" json:"asset_no"`
	AssetTag string `gorm:"size:100;column:asset_tag" json:"asset_tag"`

	Buildings   []Building `gorm:"many2many:asset_buildings;" json:"buildings"`
	BuildingIDs []uint     `gorm:"-" json:"building_ids,omitempty"`

	BrandID *uint `gorm:"NULL" json:"brand_id"`
	Brand   Brand `gorm:"foreignKey:BrandID;" json:"brand"`

	ModelName  string `gorm:"size:100;column:model" json:"model"`
	UnitNumber string `gorm:"size:50;column:unit_number" json:"unit_number"`

	SpecTransformerID *uint64         `gorm:"NULL" json:"spec_transformer_id"`
	SpecTransformer   SpecTransformer `gorm:"foreignKey:SpecTransformerID;" json:"spec_transformer"`

	SpecGeneratorID *uint64       `gorm:"NULL" json:"spec_generator_id"`
	SpecGenerator   SpecGenerator `gorm:"foreignKey:SpecGeneratorID;" json:"spec_generator"`

	SpecElectricalPanelID *uint64             `gorm:"NULL" json:"spec_electrical_panel_id"`
	SpecElectricalPanel   SpecElectricalPanel `gorm:"foreignKey:SpecElectricalPanelID" json:"spec_electrical_panel"`

	SpecFireAlarmID *uint64       `gorm:"NULL" json:"spec_fire_alarm_id"`
	SpecFireAlarm   SpecFireAlarm `gorm:"foreignKey:SpecFireAlarmID" json:"spec_fire_alarm"`

	SpecPumpID *uint64  `gorm:"NULL" json:"spec_pump_id"`
	SpecPump   SpecPump `gorm:"foreignKey:SpecPumpID;" json:"spec_pump"`

	SpecPaddleWheelID *uint64          `gorm:"NULL" json:"spec_paddle_wheel_id"`
	SpecPaddleWheel   *SpecPaddleWheel `gorm:"foreignKey:SpecPaddleWheelID;" json:"spec_paddle_wheel"`

	Status    string    `gorm:"size:50;default:'ACTIVE';column:status" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

type SpecTransformer struct {
	ID                    uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	TransformerType       string `gorm:"size:50;column:transformer_type" json:"transformer_type"`
	RatedCapacityKVA      string `gorm:"size:50;column:rated_capacity_kva" json:"rated_capacity_kva"`
	UnitOfCapacity        string `gorm:"size:20;default:'kVA';column:unit_of_capacity" json:"unit_of_capacity"`
	HVProtectionEquipment string `gorm:"size:100;column:hv_protection_equipment" json:"hv_protection_equipment"`
}

type SpecGenerator struct {
	ID               uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	EngineBrand      string   `gorm:"size:100;column:engine_brand" json:"engine_brand"`
	EngineModel      string   `gorm:"size:100;column:engine_model" json:"engine_model"`
	EnginePowerHP    *float64 `gorm:"type:numeric(10,2);column:engine_power_hp" json:"engine_power_hp,omitempty"`
	EnginePowerKW    *float64 `gorm:"type:numeric(10,2);column:engine_power_kw" json:"engine_power_kw,omitempty"`
	AlternatorBrand  string   `gorm:"size:100;column:alternator_brand" json:"alternator_brand"`
	AlternatorModel  string   `gorm:"size:100;column:alternator_model" json:"alternator_model"`
	BaseRatingKVA    *float64 `gorm:"type:numeric(10,2);column:base_rating_kva" json:"base_rating_kva,omitempty"`
	BaseRatingKW     *float64 `gorm:"type:numeric(10,2);column:base_rating_kw" json:"base_rating_kw,omitempty"`
	FuelTankCapacity *float64 `gorm:"type:numeric(10,2);column:fuel_tank_capacity" json:"fuel_tank_capacity,omitempty"`
	FuelTankUnit     string   `gorm:"size:20;default:'Litre';column:fuel_tank_unit" json:"fuel_tank_unit"`
}

type SpecElectricalPanel struct {
	ID              uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	MainType        string   `gorm:"size:100;column:main_type" json:"main_type"`
	Brand           string   `gorm:"size:100;column:brand" json:"brand"`
	RateCurrent     *float64 `gorm:"type:numeric(10,2);column:rate_current" json:"rate_current,omitempty"`
	UnitRateCurrent string   `gorm:"size:20;default:'A';column:unit_rate_current" json:"unit_rate_current"`
	InIu            string   `gorm:"size:50;column:in_iu" json:"in_iu"`
	Ue              string   `gorm:"size:50;column:ue" json:"ue"`
	Icu             string   `gorm:"size:50;column:icu" json:"icu"`
	Ics             string   `gorm:"size:50;column:ics" json:"ics"`
	Icw             string   `gorm:"size:50;column:icw" json:"icw"`
	Uimp            string   `gorm:"size:50;column:uimp" json:"uimp"`
	Hz              string   `gorm:"size:20;default:'50Hz';column:hz" json:"hz"`
	CatType         string   `gorm:"size:50;column:cat_type" json:"cat_type"`
	LSetting        string   `gorm:"size:50;column:l_setting" json:"l_setting"`
	TimeOfL         string   `gorm:"size:50;column:time_of_l" json:"time_of_l"`
	Setting         string   `gorm:"size:50;column:s_setting" json:"s_setting"`
}

type SpecFireAlarm struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemType       string `gorm:"size:50;column:system_type" json:"system_type"`
	FcpBrand         string `gorm:"size:100;column:fcp_brand" json:"fcp_brand"`
	FcpModel         string `gorm:"size:100;column:fcp_model" json:"fcp_model"`
	CabinetCount     *int   `gorm:"column:cabinet_count" json:"cabinet_count,omitempty"`
	AnnunciatorBrand string `gorm:"size:100;column:annunciator_brand" json:"annunciator_brand"`
	BatterySpecs     string `gorm:"size:255;column:battery_specs" json:"battery_specs"`
}

type SpecPump struct {
	ID                   uint     `gorm:"primaryKey;column:id" json:"id"`
	PumpApplication      string   `gorm:"size:50;column:pump_application" json:"pump_application"`
	DriveType            string   `gorm:"size:50;column:drive_type" json:"drive_type"`
	MotorBrand           string   `gorm:"size:100;column:motor_brand" json:"motor_brand"`
	MotorModel           string   `gorm:"size:100;column:motor_model" json:"motor_model"`
	PowerRating          *float64 `gorm:"type:numeric(10,2);column:power_rating" json:"power_rating,omitempty"`
	PowerUnit            string   `gorm:"size:20;column:power_unit" json:"power_unit"`
	MotorRPM             string   `gorm:"size:50;column:motor_rpm" json:"motor_rpm"`
	PumpBrand            string   `gorm:"size:100;column:pump_brand" json:"pump_brand"`
	PumpModel            string   `gorm:"size:100;column:pump_model" json:"pump_model"`
	FlowRate             string   `gorm:"size:100;column:flow_rate" json:"flow_rate"`
	FlowRateUnit         string   `gorm:"size:50;column:flow_rate_unit" json:"flow_rate_unit"`
	PressureTankCapacity string   `gorm:"size:100;column:pressure_tank_capacity" json:"pressure_tank_capacity"`
	JockeyPumpSpecs      string   `gorm:"size:255;column:jockey_pump_specs" json:"jockey_pump_specs"`
	BatterySpecs         string   `gorm:"size:100;column:battery_specs" json:"battery_specs"`
}

type SpecPaddleWheel struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PowerType       string `gorm:"size:50;column:power_type" json:"power_type"`
	MotorBrand      string `gorm:"size:100;column:motor_brand" json:"motor_brand"`
	ModelName       string `gorm:"size:100;column:model" json:"model"`
	PowerRating     string `gorm:"size:50;column:power_rating" json:"power_rating"`
	PowerUnit       string `gorm:"size:20;column:power_unit" json:"power_unit"`
	RPMEnergySource string `gorm:"size:100;column:rpm_energy_source" json:"rpm_energy_source"`
}

func main() {
	var version = "api/v1"
	fmt.Println("Starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("dsn")
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)

	} else {
		if err = db.AutoMigrate(
			&ServiceRecord{},
			&Category{},
			&Brand{},
			&Building{},
			&Asset{},
			&SpecTransformer{},
			&SpecGenerator{},
			&SpecElectricalPanel{},
			&SpecFireAlarm{},
			&SpecPump{},
			&SpecPaddleWheel{},
		); err != nil {
			log.Panic(err)
		}
	}

	if os.Getenv("production") == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if os.Getenv("production") == "true" {
		router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	} else {
		router.SetTrustedProxies(nil)
	}

	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	serviceRecordHandler := &ServiceRecordHandler{DB: db}
	serviceRecordGroup := router.Group("/" + version + "/serviceRecord")
	{
		serviceRecordGroup.GET("/", serviceRecordHandler.GetServiceRecords)
		serviceRecordGroup.GET("/:id", serviceRecordHandler.GetServiceRecord)
		serviceRecordGroup.POST("/", serviceRecordHandler.CreateServiceRecord)
		serviceRecordGroup.PUT("/:id", serviceRecordHandler.UpdateServiceRecord)
		//		serviceRecordGroup.DELETE("/:id", serviceRecordHandler.DeleteServiceRecord)
	}

	categoryHandler := &CategoryHandler{DB: db}
	categoryGroup := router.Group("/" + version + "/category")
	{
		categoryGroup.GET("/", categoryHandler.GetCategories)
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
		categoryGroup.POST("/", categoryHandler.CreateCategory)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		// categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	buildingHandler := &BuildingHandler{DB: db}
	buildingGroup := router.Group("/" + version + "/building")
	{
		buildingGroup.GET("/", buildingHandler.GetBuildings)
		buildingGroup.GET("/:building_code", buildingHandler.GetBuilding)
		buildingGroup.POST("/", buildingHandler.CreateBuilding)
		buildingGroup.PUT("/:building_code", buildingHandler.UpdateBuilding)
		//buildingGroup.DELETE("/:building_code", buildingHandler.DeleteBuilding)
	}

	brandHandler := &BrandHandler{DB: db}
	brandGroup := router.Group("/" + version + "/brand")
	{
		brandGroup.GET("/", brandHandler.GetBrands)
		brandGroup.GET("/:id", brandHandler.GetBrandByID)
		brandGroup.POST("/", brandHandler.CreateBrand)
		brandGroup.PUT("/:id", brandHandler.UpdateBrand)
		//brandGroup.DELETE("/:id", brandHandler.DeleteBrand)
	}

	assetHandler := &AssetHandler{DB: db}
	assetGroup := router.Group("/" + version + "/asset")
	{
		assetGroup.GET("/", assetHandler.GetAssets)
		assetGroup.GET("/category/:category_id", assetHandler.GetAssetsByCategory)
		assetGroup.GET("/:asset_id", assetHandler.GetAsset)
		assetGroup.POST("/", assetHandler.CreateAsset)
		assetGroup.PUT("/:asset_id", assetHandler.UpdateAsset)
		// assetGroup.DELETE("/:asset_id", assetHandler.DeleteAsset)
	}

	specTransformerHandler := &SpecTransformerHandler{DB: db}
	specTransformerGroup := router.Group("/" + version + "/specTransformer")
	{
		specTransformerGroup.GET("/", specTransformerHandler.GetSpecTransformers)
		specTransformerGroup.GET("/:spec_transformer_id", specTransformerHandler.GetSpecTransformer)
		specTransformerGroup.POST("/", specTransformerHandler.CreateSpecTransformer)
		specTransformerGroup.PUT("/:spec_transformer_id", specTransformerHandler.UpdateSpecTransformer)
		// specTransformerGroup.DELETE("/:spec_transformer_id", specTransformerHandler.DeleteSpecTransformer)
	}

	specGeneratorHandler := &SpecGeneratorHandler{DB: db}
	specGeneratorGroup := router.Group("/" + version + "/specGenerator")
	{
		specGeneratorGroup.GET("/", specGeneratorHandler.GetSpecGenerators)
		specGeneratorGroup.GET("/:spec_generator_id", specGeneratorHandler.GetSpecGenerator)
		specGeneratorGroup.POST("/", specGeneratorHandler.CreateSpecGenerator)
		specGeneratorGroup.PUT("/:spec_generator_id", specGeneratorHandler.UpdateSpecGenerator)
		// specGeneratorGroup.DELETE("/:spec_generator_id", specGeneratorHandler.DeleteSpecGenerator)
	}

	specElectricalPanelHandler := &SpecElectricalPanelHandler{DB: db}
	specElectricalPanelGroup := router.Group("/" + version + "/specElectricalPanel")
	{
		specElectricalPanelGroup.GET("/", specElectricalPanelHandler.GetSpecElectricalPanels)
		specElectricalPanelGroup.GET("/:spec_electrical_panel_id", specElectricalPanelHandler.GetSpecElectricalPanel)
		specElectricalPanelGroup.POST("/", specElectricalPanelHandler.CreateSpecElectricalPanel)
		specElectricalPanelGroup.PUT("/:spec_electrical_panel_id", specElectricalPanelHandler.UpdateSpecElectricalPanel)
		// specElectricalPanelGroup.DELETE("/:spec_electrical_panel_id", specElectricalPanelHandler.DeleteSpecElectricalPanel)
	}

	specFireAlarmHandler := &SpecFireAlarmHandler{DB: db}
	specFireAlarmGroup := router.Group("/" + version + "/specFireAlarm")
	{
		specFireAlarmGroup.GET("/", specFireAlarmHandler.GetSpecFireAlarms)
		specFireAlarmGroup.GET("/:spec_fire_alarm_id", specFireAlarmHandler.GetSpecFireAlarm)
		specFireAlarmGroup.POST("/", specFireAlarmHandler.CreateSpecFireAlarm)
		specFireAlarmGroup.PUT("/:spec_fire_alarm_id", specFireAlarmHandler.UpdateSpecFireAlarm)
		// specFireAlarmGroup.DELETE("/:spec_fire_alarm_id", specFireAlarmHandler.DeleteSpecFireAlarm)
	}

	specPumpHandler := &SpecPumpHandler{DB: db}
	specPumpGroup := router.Group("/" + version + "/specPump")
	{
		specPumpGroup.GET("/", specPumpHandler.GetSpecPumps)
		specPumpGroup.GET("/:spec_pump_id", specPumpHandler.GetSpecPump)
		specPumpGroup.POST("/", specPumpHandler.CreateSpecPump)
		specPumpGroup.PUT("/:spec_pump_id", specPumpHandler.UpdateSpecPump)
		// specPumpGroup.DELETE("/:spec_pump_id", specPumpHandler.DeleteSpecPump)
	}

	specPaddleWheelHandler := &SpecPaddleWheelHandler{DB: db}
	specPaddleWheelGroup := router.Group("/" + version + "/specPaddleWheel")
	{
		specPaddleWheelGroup.GET("/", specPaddleWheelHandler.GetSpecPaddleWheels)
		specPaddleWheelGroup.GET("/:spec_paddle_wheel_id", specPaddleWheelHandler.GetSpecPaddleWheel)
		specPaddleWheelGroup.POST("/", specPaddleWheelHandler.CreateSpecPaddleWheel)
		specPaddleWheelGroup.PUT("/:spec_paddle_wheel_id", specPaddleWheelHandler.UpdateSpecPaddleWheel)
		// specPaddleWheelGroup.DELETE("/:spec_paddle_wheel_id", specPaddleWheelHandler.DeleteSpecPaddleWheel)
	}

	if os.Getenv("production") == "true" {
		router.Run("127.0.0.1:8080") // listen and serve on

	} else {
		router.Run("localhost:8080") // listen and serve on
	}
}

func CORSMiddleware() gin.HandlerFunc {
	allowOrigins := map[string]bool{
		"http://localhost:4200": true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
