package example.authz

allow {
	input.auth.user == "alice"
}
