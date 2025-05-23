
##

```bash
strace -C -f -e trace=setsockopt  ./aileron -f _example/listener/
```

```txt
[pid 32376] setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [1], 4) = 0
[pid 32376] setsockopt(7, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
[pid 32376] setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
syscall 0xc000708860
[pid 32376] setsockopt(3, SOL_SOCKET, SO_KEEPALIVE, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_LINGER, {l_onoff=1, l_linger=123456}, 8) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_RCVBUF, [3000], 4) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_SNDTIMEO_OLD, "\2\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_RCVTIMEO_OLD, "\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 16) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_REUSEPORT, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_SOCKET, SO_SNDBUF, [4000], 4) = 0
[pid 32376] setsockopt(3, SOL_IP, IP_BIND_ADDRESS_NO_PORT, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_IP, IP_FREEBIND, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_IP, IP_TTL, [50], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_CORK, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_DEFER_ACCEPT, [10], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_KEEPCNT, [100], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_KEEPIDLE, [10], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_KEEPINTVL, [10], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_LINGER2, "\1\0\0\0\1\0\0\0", 8) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_MAXSEG, [4096], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_NODELAY, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_QUICKACK, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_SYNCNT, [100], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_USER_TIMEOUT, [1000], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_WINDOW_CLAMP, [100], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_FASTOPEN, [1], 4) = 0
[pid 32376] setsockopt(3, SOL_TCP, TCP_FASTOPEN_CONNECT, [1], 4) = -1 EINVAL (Invalid argument)
```
