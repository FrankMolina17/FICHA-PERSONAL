package models

import "gorm.io/gorm"

// TAREA (CP1): Complete los campos de Cliente según lo que muestran las pantallas.
//
// Pistas de trabajo:
//   - Guíese por el modelo Clase para los tags gorm y json.
//   - Los tests de acceptance/ compilan contra los nombres EXACTOS de los campos.
//     Mientras falten campos, `go test ./acceptance/...` mostrará errores de
//     compilación que le indican qué está faltando.
type Cliente struct {
	gorm.Model
	Nombre   string `gorm:"size:120;not null" json:"nombre"`
	Cedula   string `gorm:"size:20;not null;unique" json:"cedula"`
	Telefono string `gorm:"size:20;not null" json:"telefono"`
}
