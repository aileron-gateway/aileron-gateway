package authz

default allow = false

allow {
    input.method = "GET"
    input.api = "/allowed"
}