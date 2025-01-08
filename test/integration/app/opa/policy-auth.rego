package authz

import future.keywords.in

default allow = false

allow {
    # input.auth is a object (map[string]any of Go)
    some "admin" in input.auth["role"]
}
