# This is an example of Dockerfile.
# This Dockerfile is development use only.

##### Stage 1
# ----------*----------*----------
# /
# └─ work/
#    ├─ go.mod
#    ├─ go.sum
#    ├─ ...
#    └─ aileron
# ----------*----------*----------
FROM golang:latest AS builder
WORKDIR /work

ENV GOTOOLCHAIN=auto
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN make build

##### Stage 2
# ----------*----------*----------
# /
# └─ app/
#    ├─ _example/
#    └─ aileron
# ----------*----------*----------
FROM busybox:glibc
WORKDIR /app

RUN echo "nonroot:x:65532:" >> /etc/group
RUN echo "nonroot:x:65532:65532:nonroot:/home:/bin/false" >> /etc/passwd
USER nonroot

COPY --from=builder /work/aileron /app/aileron
COPY _example/ /app/_example/

CMD ["aileron", "-f", "_example/template/"]
