# Packaging binary

## Setup environments

### Install Golang

See <https://go.dev/doc/install>

### Install nfpm

See <https://nfpm.goreleaser.com/install/>

## Create packages

### rpm package

```sh
make rpm
```

### deb package

```sh
make deb
```

### apk package

```sh
make apk
```

## Install and Remove

### rpm package

```sh
# Install
rpm --nodeps -ivh ./aileron-0.0.0-1.x86_64.rpm
# Remove
rpm -e aileron-0.0.0-1.x86_64
```

```sh
# Install
yum install ./aileron-0.0.0-1.x86_64.rpm
# Remove
yum remove -y aileron
```

### deb package

```sh
# Install
apt install ./aileron-0.0.0-1.x86_64.deb
# Remove
apt remove --purge aileron
```

### apk package

```sh
# Install
apk add --allow-untrusted aileron-0.0.0-1.x86_64.apk
# Remove
apk del aileron
```

## References

- <https://fpm.readthedocs.io/en/v1.15.1/>
- <https://wiki.alpinelinux.org/wiki/Creating_an_Alpine_package>
- <https://nfpm.goreleaser.com/>
- <https://github.com/grafana/alloy/blob/main/tools/make/packaging.mk>
- <https://old.calculate-linux.org/main/en/openrc_manuals>
- <https://wiki.gentoo.org/wiki/OpenRC_to_systemd_Cheatsheet>
