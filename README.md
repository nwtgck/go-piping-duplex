# piping-duplex
Duplex communication over [Piping Server](https://github.com/nwtgck/piping-server)

## Installation

Get executable binaries from [GitHub Releases](https://github.com/nwtgck/go-piping-duplex/releases)

## Usage

Person A and B can communicate as follows.

Person A
```bash
piping-duplex aaa bbb
```

Person B
```bash
piping-duplex bbb aaa
```

### Example 1

```console
$ echo hello | piping-duplex aaa bbb
1
2
3
4
5
6
7
8
9
10
```

Person B
```console
$ seq 10 | piping-duplex bbb aaa | base64
aGVsbG8K
```

NOTE: `aGVsbG8K` is base64-encoded "hello\n".

### Example SSH

Here is an example to use SSH over HTTPS via Piping Server.

In server host, type as follows.
```bash
socat 'EXEC:piping-duplex aaa bbb' TCP:127.0.0.1:22
```

In client host, type as follows.
```bash
socat TCP-LISTEN:31376 'EXEC:piping-duplex bbb aaa'
````

In another terminal in client host, type as follows.

```bash
ssh -p 31376 localhost
```

### Specify Piping Server

You can specify Piping Server with `--server` or `-s` option.

```
piping-duplex -s https://piping.glitch.me aaa bbb
```

OR

specify with `$PIPING_SERVER_URL` environment variable as follows.

```bash
export PIPING_SERVER_URL=https://piping.glitch.me
```

## Help

```
Duplex communication over Piping Server

Usage:
  piping-duplex [flags]

Flags:
  -h, --help            help for piping-duplex
  -s, --server string   Piping Server URL (default "https://ppng.io")
  -c, --symmetric       use symmetric passphrase protection
  -v, --version         show version
```
