# Debian packaging assets

This folder contains assets for `.deb` packages.

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

## Install and remove with apt

**Install.**

```bash
ARCH=amd64
VERSION=v1.0.0
sudo apt install ./aileron_${VERSION}-1_${ARCH}.deb
```

**Remove.**

```bash
sudo apt remove --purge aileron
```

## Install and remove with dpkg

**Install.**

```bash
ARCH=amd64
VERSION=v1.0.0
sudo dpkg --install ./aileron_${VERSION}-1_${ARCH}.deb
```

**Remove.**

```bash
sudo dpkg --purge aileron
```
