package inventario

import (
	"errors"
	"fmt"
)

// Atributos comunes - mercaderia

type ProductoBase struct {
	id     int
	nombre string
	precio float64
	stock  int
}

type ProductoElectronico struct {
	ProductoBase
	GarantiaMeses int
}

func NuevoElectronico(id int, nombre string, precio float64, stock int, garantia int) ProductoElectronico {
	return ProductoElectronico{
		ProductoBase:  ProductoBase{id: id, nombre: nombre, precio: precio, stock: stock},
		GarantiaMeses: garantia,
	}
}

// Getters public

func (p ProductoBase) ID() int         { return p.id }
func (p ProductoBase) Nombre() string  { return p.nombre }
func (p ProductoBase) Precio() float64 { return p.precio }
func (p ProductoBase) Stock() int      { return p.stock }

// BD Simulada en memoria (Slices - mercaderia)

var CatalogoMercaderia = []ProductoElectronico{
	NuevoElectronico(1, "Laptop Asus ROG", 1250.00, 4, 12),
	NuevoElectronico(2, "Mouse Óptico Logitech", 35.50, 15, 6),
	NuevoElectronico(3, "Teclado Mecánico Corsair", 85.00, 8, 12),
}

// HistorialBaseDatos registra cada transaccion de compra

var HistorialBaseDatos []string

// RegistrarEnBaseLogs - auditoria al historial (Simulacion de DB)

func RegistrarEnBaseLogs(log string) {
	HistorialBaseDatos = append(HistorialBaseDatos, log)
}

// AgregarProducto - Expandir la mercaderia

func AgregarProducto(id int, nombre string, precio float64, stock int, garantia int) {
	nuevo := NuevoElectronico(id, nombre, precio, stock, garantia)
	CatalogoMercaderia = append(CatalogoMercaderia, nuevo)

	// Registro obligatorio - logs de bd

	logTransaccion := fmt.Sprintf("ADMIN: Producto agregado -> ID: %d | %s | Stock Inicial: %d", id, nombre, stock)
	RegistrarEnBaseLogs(logTransaccion)
}

// VerificarYReducirStock - controla el inventario

func VerificarYReducirStock(id int, cantidad int) (ProductoElectronico, error) {
	if cantidad <= 0 {
		return ProductoElectronico{}, errors.New("la cantidad debe ser mayor a cero")
	}

	for i := range CatalogoMercaderia {
		if CatalogoMercaderia[i].id == id {
			if CatalogoMercaderia[i].stock < cantidad {
				return ProductoElectronico{}, fmt.Errorf("existencias insuficientes de %s", CatalogoMercaderia[i].nombre)
			}
			// Mutacinn del estado encapsulado
			CatalogoMercaderia[i].stock -= cantidad
			return CatalogoMercaderia[i], nil
		}
	}
	return ProductoElectronico{}, errors.New("el ID del articulo no existe")
}

// MostrarMercaderia - catalogo actual

func MostrarMercaderia() {
	for _, prod := range CatalogoMercaderia {
		fmt.Printf("[%d] - %-25s | Precio: $%7.2f | Stock: %d uds | Garantia: %d meses\n",
			prod.ID(), prod.Nombre(), prod.Precio(), prod.Stock(), prod.GarantiaMeses)
	}
}

// MostrarLogsDB - operaciones registradas
func MostrarLogsDB() {
	if len(HistorialBaseDatos) == 0 {
		fmt.Println("No existen registros de transacciones en la base de datos.")
		return
	}
	for i, log := range HistorialBaseDatos {
		fmt.Printf("LOG_DB_00%d: %s\n", i+1, log)
	}
}
