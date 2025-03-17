# RPM packaging assets

This folder contains assets for `.rpm` packages.

Packages are built with [nfpm](https://github.com/goreleaser/nfpm).

## Installed files

```txt
/
├── etc/
│   ├── sysconfig/
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

## Install and remove with rpm

**Install.**

```bash
ARCH=x86_64
VERSION=v1.0.0
sudo rpm --install ./aileron-${VERSION}-1.${ARCH}.rpm
```

**Remove.**

```bash
sudo rpm --erase aileron
```

## Install and remove with yum

**Install.**

```bash
ARCH=x86_64
VERSION=v1.0.0
sudo yum install ./aileron-${VERSION}-1.${ARCH}.rpm
```

**Remove.**

```bash
sudo yum remove aileron
```

## Install and remove with dnf

**Install.**

```bash
ARCH=x86_64
VERSION=v1.0.0
sudo dnf install ./aileron-${VERSION}-1.${ARCH}.rpm
```

**Remove.**

```bash
sudo dnf erase aileron
```
