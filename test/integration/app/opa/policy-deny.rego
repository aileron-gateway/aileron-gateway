package authz

default allow = true

deny {
    input.method = "GET"
    input.api = "/denied"
    false
}