package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"myapp/models"

	"github.com/go-chi/chi/v5"
)

// getUserHandler obtiene info pública de un usuario por ID
func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener userID de la URL
		userIDStr := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}

		// Consultar solo los datos públicos (ID, Username)
		var userResp models.UserResponse // Usa la struct segura para respuestas
		err = db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).Scan(&userResp.ID, &userResp.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			} else {
				log.Printf("Error consultando datos de usuario %d: %v", userID, err)
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
			return
		}

		// Devolver los datos del usuario
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userResp)
	}
}

func GetUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener userID del contexto (inyectado por el middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok || userID == 0 {
			http.Error(w, `{"error": "No se pudo obtener ID de usuario del token"}`, http.StatusInternalServerError)
			return
		}

		// Ahora usar este userID para buscar los datos del perfil
		var userResp models.UserResponse
		err := db.QueryRow("SELECT id, username FROM users WHERE id = ?", userID).Scan(&userResp.ID, &userResp.Username)
		if err != nil {
			// ... (manejo de error: 404 si no se encuentra, 500 otros) ...
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Usuario del token no encontrado"}`, http.StatusNotFound)
			} else {
				log.Printf("Error consultando perfil para user %d: %v", userID, err)
				http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userResp)
	}
}
