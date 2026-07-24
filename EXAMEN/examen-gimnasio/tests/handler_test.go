package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/handlers"
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TestClienteHandler_CrearYObtener prueba el flujo completo:
// POST /api/v1/clientes → 201, luego GET /api/v1/clientes/{id} → 200
func TestClienteHandler_CrearYObtener(t *testing.T) {
	// Preparar repositorio en memoria y service
	repo := storage.NuevoClienteMemoria()
	svc := services.NuevoClienteService(repo)
	handler := handlers.NuevoClienteHandler(svc)

	// Crear un router con las rutas del cliente
	router := nuevoRouterCliente(handler)

	// POST /api/v1/clientes → 201
	body := `{"nombre":"María García","cedula":"1310000010","telefono":"0990000010"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code,
		"POST con datos válidos debe responder 201. Body: %s", rec.Body.String())

	// Verificar que el JSON de respuesta tiene los campos correctos
	var creado models.Cliente
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &creado))
	require.Equal(t, "María García", creado.Nombre)
	require.Equal(t, "1310000010", creado.Cedula)
	require.NotZero(t, creado.ID, "debe tener ID asignado")

	// GET /api/v1/clientes/{id} → 200
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/clientes/%d", creado.ID), nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code,
		"GET con ID válido debe responder 200")

	var obtenido models.Cliente
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &obtenido))
	require.Equal(t, creado.Nombre, obtenido.Nombre)
}

// TestClienteHandler_ObtenerInexistente prueba que GET con ID inexistente → 404
func TestClienteHandler_ObtenerInexistente(t *testing.T) {
	repo := storage.NuevoClienteMemoria()
	svc := services.NuevoClienteService(repo)
	handler := handlers.NuevoClienteHandler(svc)

	router := nuevoRouterCliente(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/clientes/99999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code,
		"GET con ID inexistente debe responder 404")
}

// TestClienteHandler_CrearDatosInvalidos prueba que POST con datos incompletos → 422
func TestClienteHandler_CrearDatosInvalidos(t *testing.T) {
	repo := storage.NuevoClienteMemoria()
	svc := services.NuevoClienteService(repo)
	handler := handlers.NuevoClienteHandler(svc)

	router := nuevoRouterCliente(handler)

	// Sin nombre (campo obligatorio)
	body := `{"cedula":"1310000011","telefono":"0990000011"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clientes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code,
		"POST sin nombre debe responder 422. Body: %s", rec.Body.String())
}

// TestInscripcionHandler_CrearInscripcion prueba POST /api/v1/inscripciones → 201
func TestInscripcionHandler_CrearInscripcion(t *testing.T) {
	// Preparar datos base
	clases := storage.NuevaClaseMemoria()
	clientes := storage.NuevoClienteMemoria()
	inscripciones := storage.NuevaInscripcionMemoria()

	clase := models.Clase{Nombre: "Spinning", PrecioUnitario: 10.0, Stock: 10, Activo: true}
	require.NoError(t, clases.Crear(&clase))
	cliente := models.Cliente{Nombre: "Pedro Ruiz", Cedula: "1310000020", Telefono: "0990000020"}
	require.NoError(t, clientes.Crear(&cliente))

	svc := services.NuevaInscripcionService(inscripciones, clases, clientes)
	handler := handlers.NuevaInscripcionHandler(svc)

	router := nuevoRouterInscripcion(handler)

	body := `{"clase_id":1,"cliente_id":1,"cantidad":2}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inscripciones", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code,
		"POST inscripción válida debe responder 201. Body: %s", rec.Body.String())

	var creado models.Inscripcion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &creado))
	require.InDelta(t, 20.0, creado.Total, 0.001, "Total: 2 x $10 = $20")
	require.Equal(t, models.EstadoPendiente, creado.Estado)
}

// --- Funciones auxiliares para crear routers de prueba ---

// nuevoRouterCliente crea un router con chi solo con las rutas de clientes
func nuevoRouterCliente(h *handlers.ClienteHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/v1/clientes", h.Crear)
	r.Get("/api/v1/clientes", h.Listar)
	r.Get("/api/v1/clientes/{id}", h.ObtenerPorID)
	return r
}

// nuevoRouterInscripcion crea un router con chi solo con las rutas de inscripciones
func nuevoRouterInscripcion(h *handlers.InscripcionHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/v1/inscripciones", h.Crear)
	r.Get("/api/v1/inscripciones", h.Listar)
	r.Get("/api/v1/inscripciones/{id}", h.ObtenerPorID)
	r.Post("/api/v1/inscripciones/{id}/retirar", h.Retirar)
	return r
}
 