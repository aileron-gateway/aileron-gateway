package example.authn

allow {
	print("22222")
	print(data.authn)
	input.authn.user == "alice"
}
