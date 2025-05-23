package example.authz

import future.keywords.if

default allow := false

allow if {
    input.header["Foo"][0] == "bar"
    input.method == "POST"
}
