package lang

func setUpEs() {
	es, _ := translator.GetTranslator("es")

	// ----- FIELDS --------

	// ----- ERRORS --------
	es.Add("internal error", "Ha ocurrido un error en el servidor.", true)
	es.Add("login error", "Usuario o contraseña incorrectos.", true)
	es.Add("unauthenticated", "Usuario no autenticado.", true)
	es.Add("not found", "No encontrado", true)

	// ----- RESPONSES --------
}
