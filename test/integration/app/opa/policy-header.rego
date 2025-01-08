package authz

import future.keywords.in

default allow = false

allow {
    # input.header is a object (map[string][]string of Go)
    some "allowed" in input.header["Test"]
}
