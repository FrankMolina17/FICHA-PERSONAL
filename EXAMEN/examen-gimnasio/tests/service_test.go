package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/services"
	"github.com/joancema/examen-gimnasio/internal/storage"
) 

// TestInscripcionService_CrearDescuento verifica que el service calcula
// correctamente el total con y sin descuento usando repositorios en memoria.
func TestInscripcionService_CrearDescuento(t *testing.T) {
	clases := storage.NuevaClaseMemoria()
	clientes := storage.NuevoClienteMemoria()
	inscripciones := storage.NuevaInscripcionMemoria()

	// Crear una clase de prueba
	clase := models.Clase{Nombre: "Spinning", PrecioUnitario: 10.0, Stock: 20, Activo: true}
	require.NoError(t, clases.Crear(&clase))

	// Crear un cliente de prueba
	cliente := models.Cliente{Nombre: "Juan Pérez", Cedula: "1310000001", Telefono: "0990000001"}
	require.NoError(t, clientes.Crear(&cliente))

	svc := services.NuevaInscripcionService(inscripciones, clases, clientes)

	// Caso 1: 3 unidades (sin descuento)
	a1 := models.Inscripcion{ClaseID: clase.ID, ClienteID: cliente.ID, Cantidad: 3}
	require.NoError(t, svc.Crear(&a1))
	require.InDelta(t, 30.0, a1.Total, 0.001, "3 x $10 = $30 sin descuento")
	require.Equal(t, models.EstadoPendiente, a1.Estado)

	// Caso 2: 5 unidades (con 10% descuento)
	a2 := models.Inscripcion{ClaseID: clase.ID, ClienteID: cliente.ID, Cantidad: 5}
	require.NoError(t, svc.Crear(&a2))
	require.InDelta(t, 45.0, a2.Total, 0.001, "5 x $10 = $50, con 10% desc = $45")

	// Verificar que el stock se descontó correctamente
	claseActual, ok := clases.ObtenerPorID(clase.ID)
	require.True(t, ok)
	require.Equal(t, uint(12), claseActual.Stock, "Stock: 20 - 3 - 5 = 12")
}

// TestInscripcionService_RetirarReposicion verifica que al retirar
// una inscripción PENDIENTE se repone el stock.
func TestInscripcionService_RetirarReposicion(t *testing.T) {
	clases := storage.NuevaClaseMemoria()
	clientes := storage.NuevoClienteMemoria()
	inscripciones := storage.NuevaInscripcionMemoria()

	clase := models.Clase{Nombre: "Yoga", PrecioUnitario: 8.0, Stock: 10, Activo: true}
	require.NoError(t, clases.Crear(&clase))

	cliente := models.Cliente{Nombre: "Ana López", Cedula: "1310000002", Telefono: "0990000002"}
	require.NoError(t, clientes.Crear(&cliente))

	svc := services.NuevaInscripcionService(inscripciones, clases, clientes)

	// Crear inscripción de 4 unidades
	a := models.Inscripcion{ClaseID: clase.ID, ClienteID: cliente.ID, Cantidad: 4}
	require.NoError(t, svc.Crear(&a))
	require.Equal(t, uint(6), getStock(clases, clase.ID), "Stock: 10 - 4 = 6")

	// Retirar la inscripción
	require.NoError(t, svc.Retirar(a.ID))
	require.Equal(t, uint(10), getStock(clases, clase.ID), "Stock: 6 + 4 = 10")

	// Verificar estado RETIRADA
	obtenido, ok := inscripciones.ObtenerPorID(a.ID)
	require.True(t, ok)
	require.Equal(t, models.EstadoRetirada, obtenido.Estado)
}

// TestInscripcionService_RetirarEstadoInvalido verifica que no se puede
// retirar una inscripción que no esté en estado PENDIENTE.
func TestInscripcionService_RetirarEstadoInvalido(t *testing.T) {
	clases := storage.NuevaClaseMemoria()
	clientes := storage.NuevoClienteMemoria()
	inscripciones := storage.NuevaInscripcionMemoria()

	clase := models.Clase{Nombre: "Crossfit", PrecioUnitario: 12.0, Stock: 5, Activo: true}
	require.NoError(t, clases.Crear(&clase))

	cliente := models.Cliente{Nombre: "Carlos Ruiz", Cedula: "1310000003", Telefono: "0990000003"}
	require.NoError(t, clientes.Crear(&cliente))

	svc := services.NuevaInscripcionService(inscripciones, clases, clientes)

	// Crear inscripción y cambiar estado a ASISTIDA
	a := models.Inscripcion{ClaseID: clase.ID, ClienteID: cliente.ID, Cantidad: 1}
	require.NoError(t, svc.Crear(&a))

	a.Estado = models.EstadoAsistida
	require.NoError(t, inscripciones.Actualizar(&a))

	// Intentar retirar → debe fallar
	err := svc.Retirar(a.ID)
	require.ErrorIs(t, err, services.ErrEstadoInvalido,
		"retirar una inscripción ASISTIDA debe devolver ErrEstadoInvalido")
}


func getStock(clases *storage.ClaseMemoria, id uint) uint {
	c, ok := clases.ObtenerPorID(id)
	if !ok {
		return 0
	}
	return c.Stock
}
