package lang

func setUpEn() {
	en, _ := translator.GetTranslator("en")

	// ERRORS ---------------
	en.Add("internal error", "An unexpected error has occurred.", true)
	en.Add("login error", "Incorrect username or password.", true)
	en.Add("unauthenticated", "Unauthenticated user.", true)
}
