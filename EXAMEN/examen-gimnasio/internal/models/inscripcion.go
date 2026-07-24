package models

import "gorm.io/gorm"

// TAREA (CP1): Complete los campos de Inscripcion según lo que muestran las pantallas.
//
// Pistas de trabajo:
//   - Un Inscripcion referencia a una Clase y a un Cliente (claves foráneas).
//   - Recuerde el campo de estado (use las constantes de estados.go) y el total.
//   - Los tests de acceptance/ compilan contra los nombres EXACTOS de los campos.
type Inscripcion struct {
	gorm.Model
	ClaseID   uint    `gorm:"not null" json:"clase_id"`
	ClienteID uint    `gorm:"not null" json:"cliente_id"`
	Cantidad  uint    `gorm:"not null" json:"cantidad"`
	Estado    string  `gorm:"size:20;not null;default:PENDIENTE" json:"estado"`
	Total     float64 `gorm:"not null" json:"total"`
}
