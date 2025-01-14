package app

import (
	"UrfuNavigator-backend/internal/database"
	middleware "UrfuNavigator-backend/internal/middlewares"
	"UrfuNavigator-backend/internal/objstore"
	"UrfuNavigator-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type API struct {
	Port           string
	Store          database.Store
	ObjectStore    objstore.ObjectStore
	Services       services.Services
	AllowedOrigins string
}

func NewAPI(port string, store database.Store, objectStore objstore.ObjectStore, services services.Services, allowedOrigins string) *API {
	return &API{
		Port:           port,
		Store:          store,
		ObjectStore:    objectStore,
		Services:       services,
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

	app.Get("/icon", middleware.JWTProtected, s.GetIconHandler)
	app.Get("/icons", middleware.JWTProtected, s.GetAllIconsHandler)
	app.Get("/institute", middleware.JWTProtected, s.GetInstituteHandler)
	app.Get("/institutes", middleware.JWTProtected, s.GetAllInstitutesHandler)
	app.Get("/floor", middleware.JWTProtected, s.GetFloorHandler)
	app.Get("/floors", middleware.JWTProtected, s.GetAllFloorsHandler)
	app.Get("/graph", middleware.JWTProtected, s.GetGraphHandler)
	app.Get("/graphs", middleware.JWTProtected, s.GetAllGraphsHandler)
	app.Get("/stair", middleware.JWTProtected, s.GetStairHandler)
	app.Get("/stairs", middleware.JWTProtected, s.GetAllStairsHandler)
	app.Get("/user", middleware.JWTProtected, s.GetUserHandler)
	app.Get("/users", middleware.JWTProtected, s.GetAllUsersHandler)
	app.Get("/login", s.LoginHandler)

	app.Post("/icon", middleware.JWTProtected, s.PostIconHandler)
	app.Post("/institute", middleware.JWTProtected, s.PostInstituteHandler)
	app.Post("/floor", middleware.JWTProtected, s.PostFloorFromFileHandler)
	app.Post("/register", s.RegisterHandler)

	app.Put("/institute", middleware.JWTProtected, s.PutInstituteHandler)
	app.Put("/floor", middleware.JWTProtected, s.PutFloorHandler)
	app.Put("/graph", middleware.JWTProtected, s.PutGraphHandler)
	app.Put("/stair", middleware.JWTProtected, s.PutStairHandler)

	app.Delete("/icon", middleware.JWTProtected, s.DeleteIconHandler)
	app.Delete("/institute", middleware.JWTProtected, s.DeleteInstituteHandler)
	app.Delete("/floor", middleware.JWTProtected, s.DeleteFloorHandler)
	app.Delete("/graph", middleware.JWTProtected, s.DeleteGraphHandler)
	app.Delete("/stair", middleware.JWTProtected, s.DeleteStairHandler)

	return app.Listen(s.Port)
}
