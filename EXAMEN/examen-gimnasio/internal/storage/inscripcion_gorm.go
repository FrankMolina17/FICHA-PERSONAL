package storage

import (
	"gorm.io/gorm"

	"github.com/joancema/examen-gimnasio/internal/models"
)

// TAREA (CP2): Implemente InscripcionGORM contra la interfaz InscripcionRepository.
//
// Reglas:
//   - NO cambie el nombre del tipo, del constructor ni las firmas de los métodos.
//   - Guíese por ClaseGORM: es el mismo patrón con una entidad distinta.
//   - Recuerde: aquí NO va lógica de negocio. Solo persistencia.
type InscripcionGORM struct {
	db *gorm.DB
}

func NuevaInscripcionGORM(db *gorm.DB) *InscripcionGORM {
	return &InscripcionGORM{db: db}
}

func (r *InscripcionGORM) Crear(a *models.Inscripcion) error {
	return r.db.Create(a).Error
}

func (r *InscripcionGORM) ObtenerPorID(id uint) (models.Inscripcion, bool) {
	var a models.Inscripcion
	if err := r.db.First(&a, id).Error; err != nil {
		return models.Inscripcion{}, false
	}
	return a, true
}

func (r *InscripcionGORM) Listar() ([]models.Inscripcion, error) {
	var lista []models.Inscripcion
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *InscripcionGORM) Actualizar(a *models.Inscripcion) error {
	return r.db.Save(a).Error
}
 