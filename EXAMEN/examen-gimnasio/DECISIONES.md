# DECISIONES (CP1) — máximo 10 líneas

1. Cliente: pantalla 02 muestra formulario con nombre, cédula y teléfono para identificar al socio.
2. Clase: pantalla 01 lista clases con nombre, precio unitario, stock y estado activo/inactivo.
3. Inscripcion: pantalla 02 referencia cliente y clase; pantalla 03 muestra cantidad, estado y total.
4. Inscripcion tiene FK ClaseID y ClienteID porque une un cliente a una clase con cantidad y estado.
5. Estado se define como constante (PENDIENTE/ASISTIDA/RETIRada) para controlar el flujo.
6. Total = cantidad × precio unitario, con 10% descuento si cantidad ≥ 5 (regla R3).
7. Stock se descuenta al crear y se repone al retirar (reglas R2 y R5).
8. Cliente requiere nombre, cédula (única) y teléfono obligatorios.


//cliente.go
//inscripcion.go
//cliente_gorm.go
//cliente_service.go
//cliente_handler.go

Commit: CP1: modelos completos + vertical Cliente (repo/service/handler) + DECISIONES.md
Commit: CP2: InscripcionGORM + InscripcionService con R1-R5
Commit: CP3: ClienteService con R1-R5 - Expuesto y aprobado



