package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
)

// TAREA (CP3): Implemente InscripcionHandler.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Mapeo de errores de dominio a status codes (los tests lo verifican):
//       ErrDatosInvalidos     -> 422 Unprocessable Entity
//       ErrReferenciaInvalida -> 422 Unprocessable Entity
//       ErrStockInsuficiente  -> 409 Conflict
//       ErrEstadoInvalido     -> 409 Conflict 
//       ErrNoEncontrado       -> 404 Not Found
//       cualquier otro error  -> 500 Internal Server Error 
type InscripcionHandler struct {
	servicio *services.InscripcionService
}

func NuevaInscripcionHandler(s *services.InscripcionService) *InscripcionHandler { 
	return &InscripcionHandler{servicio: s}
}

func (h *InscripcionHandler) Crear(w http.ResponseWriter, r *http.Request) {
	var inscripcion models.Inscripcion
	if err := json.NewDecoder(r.Body).Decode(&inscripcion); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := h.servicio.Crear(&inscripcion); err != nil {
		switch {
		case errors.Is(err, services.ErrDatosInvalidos), errors.Is(err, services.ErrReferenciaInvalida):
			RespondError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, services.ErrStockInsuficiente):
			RespondError(w, http.StatusConflict, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	RespondJSON(w, http.StatusCreated, inscripcion)
}

func (h *InscripcionHandler) Listar(w http.ResponseWriter, r *http.Request) {
	lista, err := h.servicio.Listar()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lista)
}

func (h *InscripcionHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	inscripcion, err := h.servicio.ObtenerPorID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNoEncontrado):
			RespondError(w, http.StatusNotFound, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	RespondJSON(w, http.StatusOK, inscripcion)
}

func (h *InscripcionHandler) Retirar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.servicio.Retirar(uint(id)); err != nil {
		switch {
		case errors.Is(err, services.ErrNoEncontrado):
			RespondError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrEstadoInvalido):
			RespondError(w, http.StatusConflict, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "inscripción retirada"})
}
