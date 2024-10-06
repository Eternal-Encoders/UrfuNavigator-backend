package app

import (
	"UrfuNavigator-backend/internal/database"
	"UrfuNavigator-backend/internal/objstore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	Port           string
	Store          database.Store
	ObjectStore    objstore.ObjectStore
	AllowedOrigins string
}

func NewAPI(port string, store database.Store, objectStore objstore.ObjectStore, allowedOrigins string) *API {
	return &API{
		Port:           port,
		Store:          store,
		ObjectStore:    objectStore,
		AllowedOrigins: allowedOrigins,
	}
}

func (s *API) Run() error {
	app := fiber.New()

	cors := cors.New(cors.Config{
		AllowOrigins: s.AllowedOrigins,
	})

	// cfg := swagger.Config{
	// 	BasePath: "/",
	// 	Path:     "swagger",
	// 	Title:    "Swagger API Docs",
	// }

	// app.Use(swagger.New(cfg))
	app.Use(cors)

	app.Get("/icon", s.GetIconHandler)
	app.Get("/icons", s.GetAllIconsHandler)
	app.Get("/institute", s.GetInstituteHandler)
	app.Get("/institutes", s.GetAllInstitutesHandler)
	app.Post("/icon", s.PostIconHandler)
	//app.Post()

	return app.Listen(s.Port)
}
