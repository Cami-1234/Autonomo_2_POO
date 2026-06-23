package usuarios

// Propiedades base de un usuario

type Usuario struct {
	nombre string
	rol    string // Admin o cliente
}

// NuevoUsuario
func NuevoUsuario(nombre string, rol string) Usuario {
	return Usuario{nombre: nombre, rol: rol}
}

// Getters - Info del usuaro

func (u Usuario) Nombre() string { return u.nombre }
func (u Usuario) Rol() string    { return u.rol }
