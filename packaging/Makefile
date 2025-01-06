PKG_DIR ?= $(CURDIR)/pkgs/

.PHONY: rpm
rpm:
	mkdir -p $(PKG_DIR)
	# PKG_ARCH=i386    BIN_ARCH=386     nfpm package -p rpm -t $(PKG_DIR)
	PKG_ARCH=x86_64  BIN_ARCH=amd64   nfpm package -p rpm -t $(PKG_DIR)
	# PKG_ARCH=armv7   BIN_ARCH=arm     nfpm package -p rpm -t $(PKG_DIR)
	# PKG_ARCH=aarch64 BIN_ARCH=arm64   nfpm package -p rpm -t $(PKG_DIR)
	# PKG_ARCH=ppc64le BIN_ARCH=ppc64le nfpm package -p rpm -t $(PKG_DIR)
	# PKG_ARCH=riscv64 BIN_ARCH=riscv64 nfpm package -p rpm -t $(PKG_DIR)
	# PKG_ARCH=s390x   BIN_ARCH=s390x   nfpm package -p rpm -t $(PKG_DIR)

.PHONY: deb
deb:
	mkdir -p $(PKG_DIR)
	# PKG_ARCH=i386    BIN_ARCH=386     nfpm package -p deb -t $(PKG_DIR)
	PKG_ARCH=amd64   BIN_ARCH=amd64   nfpm package -p deb -t $(PKG_DIR)
	# PKG_ARCH=armhf   BIN_ARCH=arm     nfpm package -p deb -t $(PKG_DIR)
	# PKG_ARCH=arm64   BIN_ARCH=arm64   nfpm package -p deb -t $(PKG_DIR)
	# PKG_ARCH=ppc64el BIN_ARCH=ppc64le nfpm package -p deb -t $(PKG_DIR)
	# PKG_ARCH=riscv64 BIN_ARCH=riscv64 nfpm package -p deb -t $(PKG_DIR)
	# PKG_ARCH=s390x   BIN_ARCH=s390x   nfpm package -p deb -t $(PKG_DIR)

.PHONY: apk
apk:
	mkdir -p $(PKG_DIR)
	# PKG_ARCH=i386    BIN_ARCH=386     nfpm package -p apk -t $(PKG_DIR)
	PKG_ARCH=x86_64  BIN_ARCH=amd64   nfpm package -p apk -t $(PKG_DIR)
	# PKG_ARCH=armv7   BIN_ARCH=arm     nfpm package -p apk -t $(PKG_DIR)
	# PKG_ARCH=aarch64 BIN_ARCH=arm64   nfpm package -p apk -t $(PKG_DIR)
	# PKG_ARCH=ppc6le  BIN_ARCH=ppc64le nfpm package -p apk -t $(PKG_DIR)
	# PKG_ARCH=riscv64 BIN_ARCH=riscv64 nfpm package -p apk -t $(PKG_DIR)
	# PKG_ARCH=s390x   BIN_ARCH=s390x   nfpm package -p apk -t $(PKG_DIR)

.PHONY: archlinux
archlinux:
	mkdir -p $(PKG_DIR)
	PKG_ARCH=x86_64  BIN_ARCH=amd64   nfpm package -p archlinux -t $(PKG_DIR)
	# PKG_ARCH=armv7   BIN_ARCH=arm     nfpm package -p archlinux -t $(PKG_DIR)
	# PKG_ARCH=aarch64 BIN_ARCH=arm64   nfpm package -p archlinux -t $(PKG_DIR)
	# PKG_ARCH=riscv64 BIN_ARCH=riscv64 nfpm package -p archlinux -t $(PKG_DIR)
