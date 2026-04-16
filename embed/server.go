package embed

import (
	"context"
	"embed"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Recipe struct {
	Name        string
	Description string
	Time        time.Duration
}

//go:embed app/*
var appFiles embed.FS

func Execute(ctx context.Context) {
	// configure template rendering
	tmplFS := echo.MustSubFS(appFiles, "app/templates")
	renderer := NewTemplateRenderer(tmplFS)

	// create a new echo instance
	e := echo.New()
	e.Renderer = renderer
	e.Filesystem = echo.MustSubFS(appFiles, "app/assets")

	// 1. Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(ApexLoggerMiddleware)
	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogRequestID: true,
			LogHost:      true,
			LogURI:       true,
			LogStatus:    true,
			LogMethod:    true,
			LogLatency:   true,
			LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
				log.WithFields(log.Fields{
					"status":     v.Status,
					"host":       v.Host,
					"method":     v.Method,
					"uri":        v.URI,
					"request_id": v.RequestID,
					"latency":    v.Latency,
				}).Info("request")
				return nil
			},
		}))
	e.Use(middleware.Recover())

	// 2. Mock Data
	db := map[string][]Recipe{
		"drinks": {
			{Name: "Old Fashioned", Description: "A classic bourbon cocktail.", Time: 5 * time.Minute},
			{Name: "Spiced Chai", Description: "Warm tea with star anise.", Time: 15 * time.Minute},
		},
		"slow-cooker": {
			{Name: "Pulled Pork", Description: "Tender shoulder with cider vinegar.", Time: 510 * time.Minute},
			{Name: "Beef Stew", Description: "Hearty root vegetables and chuck.", Time: 6 * time.Hour},
		},
		"bbq": {
			{Name: "Smoked Ribs", Description: "Dry rubbed with hickory smoke.", Time: 5 * time.Hour},
			{Name: "Grilled Corn", Description: "Street-style with lime and cotija.", Time: 20 * time.Minute},
		},
	}

	// 3. Static Routes
	e.File("/favicon.ico", "images/favicon.ico")
	e.Static("/images", "images")
	e.Static("/js", "js")

	// 4. Dynamic Routes
	e.GET("/", func(c *echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/category/drinks")
	})
	e.GET("/category/:name", func(c *echo.Context) error {
		category := c.Param("name")
		recipes, ok := db[category]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}

		return c.Render(http.StatusOK, pageHome, map[string]any{
			"Category": strings.ReplaceAll(category, "-", " "),
			"Recipes":  recipes,
		})
	})

	err := e.Start(":8080")
	if err != nil {
		log.WithError(err).Error("server exited with errors")
	}
}
