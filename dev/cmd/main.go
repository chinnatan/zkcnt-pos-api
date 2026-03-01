package main

import (
	"log"
	"net/http"
	"zkcnt-pos-api/internal/auth"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// สร้าง global middleware สำหรับ handle CORS
		e.Router.Bind(apis.CORS(apis.CORSConfig{
			AllowOrigins:     []string{"http://localhost:8080", "https://dev-pos.zkcnt.com"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		}))

		// สร้าง global middleware สำหรับ handle error
		e.Router.BindFunc(func(e *core.RequestEvent) error {
			err := e.Next()
			if err == nil {
				return nil
			}

			if apiErr, ok := err.(*router.ApiError); ok {
				switch apiErr.Status {
				case http.StatusBadRequest:
					if emailErr, ok := apiErr.Data["email"]; ok {
						return e.JSON(http.StatusBadRequest, map[string]string{"message": emailErr.(string)})
					}
					return e.JSON(http.StatusBadRequest, map[string]string{"message": apiErr.Message})
				case http.StatusInternalServerError:
					app.Logger().Error(apiErr.Message, "error", apiErr.Data)
					if _, ok := apiErr.Data["email"]; ok {
						return e.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่สามารถสร้างบัญชีได้ เนื่องจากมี email นี้อยู่ในระบบแล้ว"})
					}
					return apiErr
				}
			}

			return err
		})

		g := e.Router.Group("/api/v1")

		auth.RegisterRoutes(app, g.Group("/auth"))
		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
