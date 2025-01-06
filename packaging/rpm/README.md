# RPM packaging assets

This folder contains assets for debian packages `.rpm`.

Packages are built with [fpm](https://fpm.readthedocs.io/).

To know the details about `.rpm` packages, check the following references.

- [Fedora Package Maintainers](https://docs.fedoraproject.org/en-US/package-maintainers/)
- [Scriptlets](https://docs.fedoraproject.org/en-US/packaging-guidelines/Scriptlets/)
- [Maximum RPM](https://ftp.osuosl.org/pub/rpm/max-rpm/)
- [RPM Packaging Guide](https://rpm-packaging-guide.github.io/)

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
