package main

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"

	"net/http"

	"github.com/djumpen/go-rest-admin/api"
	"github.com/djumpen/go-rest-admin/config"
	"github.com/djumpen/go-rest-admin/middleware/handler"
	"github.com/djumpen/go-rest-admin/util/validation"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetConfig()

	if cfg.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	gormDB, err := gorm.Open("postgres", config.GetPostgresConnection())
	if err != nil {
		log.Fatal(err)
	}

	// ---- Apply migrations -----
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(gormDB.DB(), "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}

	// Upgrade gin validator
	binding.Validator = new(validation.DefaultValidator)

	responder := api.NewResponder()

	r := gin.Default()
	r.RedirectTrailingSlash = true

	r.Use(
		cors.New(api.GetCorsConfig()),
		handler.ErrorHandler(responder, api.RegisteredRequests()),
	)

	// ------------------------ Register Resources ------------------------

	// mainSettingsSvc := services.NewMainSettings(mainSettingsSt, gormDB)

	commonRes := api.NewCommonResource(responder)
	// mainSettingsRes := api.NewMainSettingsResource(mainSettingsSvc, responder)
	// imageRes := api.NewImageResource(responder)

	// Routes

	// r.GET("/main-settings", mainSettingsRes.Read)
	// r.PATCH("/main-settings", mainSettingsRes.Update)
	// r.POST("/image/", imageRes.Upload)

	r.GET("/sanity", commonRes.Sanity)
	r.NoRoute(commonRes.NotFound)

	// --------------------------------------------------------------------

	if gin.Mode() == gin.DebugMode {
		api.PrintRegisteredRequestsReminder()
	}

	useSSL := cfg.CertFile != "" && cfg.KeyFile != ""
	address := fmt.Sprintf(":%d", cfg.Port)

	if useSSL {
		if err := http.ListenAndServeTLS(address, cfg.CertFile, cfg.KeyFile, r); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	} else {
		log.Printf("Listening on port %d\n", cfg.Port)
		if err := http.ListenAndServe(address, r); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}
}
