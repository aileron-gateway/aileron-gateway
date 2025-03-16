# Arch Linux packaging assets

This folder contains assets for `.pkg.tar.zst` packages.

Packages are built with [nfpm](https://github.com/goreleaser/nfpm).

## Installed files

```txt
/
├── etc/
│   ├── default/
│   │   └── aileron.env
│   └── aileron/
│       └── config.yaml
├── usr/
│   ├── bin/
│   │   └── aileron
│   └── lib/
│       └── systemd/
│           └── system/
│               └── aileron.service
└── var/
    └── lib/
        └── aileron/
```
