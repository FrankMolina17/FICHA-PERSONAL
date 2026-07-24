package services

import (
	"github.com/joancema/examen-gimnasio/internal/models"
	"github.com/joancema/examen-gimnasio/internal/storage"
)

// TAREA (CP2): Implemente InscripcionService con las 5 reglas de negocio.
//
// Las reglas están A LA VISTA en las pantallas (carpeta pantallas/) y los
// tests de acceptance/reglas_test.go las verifican una por una. Devuelva los
// errores de dominio de errores.go: los tests los comprueban con errors.Is.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Observe que el service recibe TRES repositories: necesita consultar
//     Clase y Cliente para validar, y actualizar Clase al retirar.
type InscripcionService struct {
	inscripciones   storage.InscripcionRepository
	clases storage.ClaseRepository
	clientes     storage.ClienteRepository
}

func NuevaInscripcionService(
	inscripciones storage.InscripcionRepository,
	clases storage.ClaseRepository,
	clientes storage.ClienteRepository,
) *InscripcionService {
	return &InscripcionService{
		inscripciones:   inscripciones,
		clases: clases,
		clientes:     clientes,
	}
}

// Crear registra una nueva inscripción aplicando R1, R2 y R3.
// R1: la clase debe existir y estar activa; el cliente debe existir.
// R2: la cantidad no puede superar el stock disponible.
// R3: total = cantidad × precio unitario, con 10% descuento si cantidad >= 5.
func (s *InscripcionService) Crear(a *models.Inscripcion) error {
	clase, ok := s.clases.ObtenerPorID(a.ClaseID)
	if !ok || !clase.Activo {
		return ErrReferenciaInvalida
	}

	_, ok = s.clientes.ObtenerPorID(a.ClienteID)
	if !ok {
		return ErrReferenciaInvalida
	}

	if a.Cantidad > clase.Stock {
		return ErrStockInsuficiente
	}

	total := float64(a.Cantidad) * clase.PrecioUnitario
	if a.Cantidad >= 5 {
		total = total * 0.90
	}
	a.Total = total
	a.Estado = models.EstadoPendiente

	if err := s.inscripciones.Crear(a); err != nil {
		return err
	}

	clase.Stock -= a.Cantidad
	return s.clases.Actualizar(&clase)
}

func (s *InscripcionService) ObtenerPorID(id uint) (models.Inscripcion, error) {
	a, ok := s.inscripciones.ObtenerPorID(id)
	if !ok {
		return models.Inscripcion{}, ErrNoEncontrado
	}
	return a, nil
}

func (s *InscripcionService) Listar() ([]models.Inscripcion, error) {
	return s.inscripciones.Listar()
}

// Retirar cancela una inscripción aplicando R4 y R5.
// R4: solo se puede retirar una inscripción en estado PENDIENTE.
// R5: al retirar, la cantidad se repone al stock de la clase.
func (s *InscripcionService) Retirar(id uint) error {
	a, ok := s.inscripciones.ObtenerPorID(id)
	if !ok {
		return ErrNoEncontrado
	}

	if a.Estado != models.EstadoPendiente {
		return ErrEstadoInvalido
	}

	a.Estado = models.EstadoRetirada
	if err := s.inscripciones.Actualizar(&a); err != nil {
		return err
	}

	clase, ok := s.clases.ObtenerPorID(a.ClaseID)
	if !ok {
		return ErrReferenciaInvalida
	}
	clase.Stock += a.Cantidad
	return s.clases.Actualizar(&clase)
}
