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
	app.Get("/floor", s.GetFloorHandler)
	app.Get("/floors", s.GetAllFloorsHandler)
	app.Get("/graph", s.GetGraphHandler)
	app.Get("/graphs", s.GetAllGraphsHandler)
	app.Get("/stair", s.GetStairHandler)
	app.Get("/stairs", s.GetAllStairsHandler)
	app.Post("/icon", s.PostIconHandler)
	app.Post("/institute", s.PostInstituteHandler)
	app.Post("/floor", s.PostFloorFromFileHandler)
	app.Delete("/icon", s.DeleteIconHandler)
	app.Delete("/institute", s.DeleteInstituteHandler)
	app.Delete("/floor", s.DeleteFloorHandler)
	app.Delete("/graph", s.DeleteGraphHandler)
	app.Delete("/stair", s.DeleteStairHandler)
	app.Put("/institute", s.PutInstituteHandler)

	return app.Listen(s.Port)
}
