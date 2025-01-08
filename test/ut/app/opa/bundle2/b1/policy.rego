package example.authn

allow {
	print("111111")
	print(input)
	print(data.foo) # From Store 
	input.authn.user == "alice"
}
