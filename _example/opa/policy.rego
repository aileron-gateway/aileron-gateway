package example.authz

import future.keywords.if
import future.keywords.in

default allow := false


allow if {
    "readonly" in input.auth.attrs.scope
    input.method == "GET"
}

allow if {
    "update" in input.auth.attrs.scope
    input.method == "POST"
}
