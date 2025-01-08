package authz

default allow = false

allow {
    input.method = "GET"
    input.api = "/allowed"
    cause_error
}

deny {
    input.method = "GET"
    input.api = "/denied"
}