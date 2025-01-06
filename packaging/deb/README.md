# Debian packaging assets

This folder contains assets for debian packages `.deb`.

Packages are built with [fpm](https://fpm.readthedocs.io/).

To know the details about `.dep` packages, check the following references.

- [6. Package maintainer scripts and installation procedure](https://www.debian.org/doc/debian-policy/ch-maintainerscripts.html)
- [9. Maintainer script flowcharts](https://www.debian.org/doc/debian-policy/ap-flowcharts.html)
- [Manpages of dpkg-dev in Debian testing](https://manpages.debian.org/testing/dpkg-dev/)

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
