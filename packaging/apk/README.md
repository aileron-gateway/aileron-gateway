# Alpine packaging assets

This folder contains assets for `.apk` packages.

Packages are built with [nfpm](https://github.com/goreleaser/nfpm).

## Installed files

```txt
/
├── etc/
│   ├── init.d/
│   │   └── aileron
│   ├── default/
│   │   └── aileron.env
│   └── aileron/
│       └── config.yaml
├── usr/
│   └── bin/
│       └── aileron
└── var/
    └── lib/
        └── aileron/
```

## Install and remove with apk

**Install.**

```bash
ARCH=x86_64
VERSION=v1.0.0
apk add --allow-untrusted ./aileron_${VERSION}-r1_${ARCH}.apk 
```

**Remove.**

```bash
apk del --purge aileron
```
