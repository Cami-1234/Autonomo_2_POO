package main

import (
	"Shopez/inventario"
	"Shopez/usuarios"
	"fmt"
)

type EvaluadorDescuento interface {
	Procesar(monto float64) float64
}

type DescuentoUIDE struct{}

func (d DescuentoUIDE) Procesar(monto float64) float64 { return monto * 0.90 } // 10%

type TarifaNormal struct{}

func (tn TarifaNormal) Procesar(monto float64) float64 { return monto }

func main() {

	// Inicio de Usuarios

	adminUser := usuarios.NuevoUsuario("Administrador", "Admin")
	clienteUser := usuarios.NuevoUsuario("Consumidor Final", "Cliente")

	var opcionMenu, sessionRol int
	fmt.Println("=======================================================")
	fmt.Println("                BIENVENIDO SHOPEZ                      ")
	fmt.Println("=======================================================")
	fmt.Printf("Seleccione Perfil de Acceso:\n[1] %s (Rol: %s)\n[2] %s (Rol: %s)\n",
		adminUser.Nombre(), adminUser.Rol(), clienteUser.Nombre(), clienteUser.Rol())
	fmt.Print("Opción> ")
	fmt.Scanln(&sessionRol)

	var usuarioActivo usuarios.Usuario
	if sessionRol == 1 {
		usuarioActivo = adminUser
	} else {
		usuarioActivo = clienteUser
	}

	for {
		fmt.Printf("\n--- MENÚ DE %s ---\n", usuarioActivo.Rol())
		if usuarioActivo.Rol() == "Admin" {
			fmt.Println("[1] Agregar Nueva Mercaderia ")
			fmt.Println("[2] Visualizar Logs e Historial ")
			fmt.Println("[3] Salir del Sistema")
		} else {
			fmt.Println("[1] Ver Catalogo de Mercaderia Disponibles")
			fmt.Println("[2] Registrar Compra de Artículo")
			fmt.Println("[3] Salir del Sistema")
		}
		fmt.Print("Seleccione una accion> ")
		fmt.Scanln(&opcionMenu)

		if opcionMenu == 3 {
			fmt.Println("Cerrando sesion en ShopEZ.")
			break
		}

		// LOGICA - ADMINISTRADOR
		if usuarioActivo.Rol() == "Admin" {
			switch opcionMenu {
			case 1:
				var id, stock, garantia int
				var nombre string
				var precio float64
				fmt.Print("Ingrese ID unico: ")
				fmt.Scanln(&id)
				fmt.Print("Nombre del articulo: ")
				fmt.Scanln(&nombre)
				fmt.Print("Precio Unitario: ")
				fmt.Scanln(&precio)
				fmt.Print("Stock Inicial: ")
				fmt.Scanln(&stock)
				fmt.Print("Meses de Garantia: ")
				fmt.Scanln(&garantia)

				inventario.AgregarProducto(id, nombre, precio, stock, garantia)
				fmt.Println(" Mercaderia registrada en bodega y auditoria ")
			case 2:
				fmt.Println("\n --- REGISTROS EN LA BASE DE DATOS --- ")
				inventario.MostrarLogsDB()
			}
		} else {

			// LOGICA - CLIENTE

			switch opcionMenu {
			case 1:
				fmt.Println("\n--- CATaLOGO DE MERCADERiA - SHOPEZ ---")
				inventario.MostrarMercaderia()
			case 2:
				var id, cant int
				var cupon string
				inventario.MostrarMercaderia()
				fmt.Print("\nID del producto a comprar: ")
				fmt.Scanln(&id)
				fmt.Print("Cantidad requerida: ")
				fmt.Scanln(&cant)

				prod, err := inventario.VerificarYReducirStock(id, cant)
				if err != nil {
					fmt.Printf("Error Operacional: %v\n", err)
				} else {
					subtotal := prod.Precio() * float64(cant)
					iva := subtotal * 0.15
					total := subtotal + iva

					fmt.Print("Ingrese cupon promocional (o 'NO'): ")
					fmt.Scanln(&cupon)

					var evaluador EvaluadorDescuento
					if cupon == "UIDE10" {
						evaluador = DescuentoUIDE{}
						fmt.Println("¡Aplicado descuento estudiantil del 10%!")
					} else {
						evaluador = TarifaNormal{}
					}

					totalFinal := evaluador.Procesar(total)
					fmt.Printf("COMPRA REALIZADA CON ÉXITO. Total Neto a Pagar: $%7.2f\n", totalFinal)

					// REGISTRO - COMPRA EN LA BASE DE DATA (LOGS)

					logCompra := fmt.Sprintf("CLIENTE: Compra efectuada -> %d unidades de %s | Pago Neto: $%7.2f", cant, prod.Nombre(), totalFinal)
					inventario.RegistrarEnBaseLogs(logCompra)
				}
			}
		}
	}
}
