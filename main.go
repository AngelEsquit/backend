package main

import (
	"log"
	"net/http"

	"myapp/handlers"
	"myapp/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db, err := setupDatabase("./users.db")

	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(configureCORS()) // CORS primero

	// --- Rutas Públicas ---
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.PostRegisterHandler(db)) // Mover register aquí
		r.Post("/login", handlers.PostLoginHandler(db))       // Mover login aquí
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { /* ... */ })

	// --- Rutas Protegidas ---
	r.Group(func(r chi.Router) {
		// Aplicar middleware JWT a este grupo
		r.Use(utils.JwtAuthMiddleware(db))

		// Rutas que requieren token válido
		r.Post("/auth/logout", handlers.PostLogoutHandler(db))      // Mover logout aquí
		r.Get("/users/profile", handlers.GetUserProfileHandler(db)) // Nueva ruta para perfil
		// La ruta /users/{userID} podría seguir siendo pública o protegerse también
		// r.Get("/users/{userID}", getUserHandler(db)) // Ejemplo si se protege
	})

	// (Opcional: Mantener /users/{userID} pública si se desea)
	r.Get("/users/{userID}", handlers.GetUserHandler(db))

	port := ":3000"
	log.Printf("Servidor escuchando en puerto %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
