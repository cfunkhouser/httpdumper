# HTTP Request Dumper

`httpdumper` accepts any request and dumps it as text as a response.

## The Server

You can run the server directly from this repository:

```console
$ go run ./cmd/httpdumper
INFO[0000] Server starting                               address=":8080"
```

You can install it directly, and run it:

```console
$ go install github.com/cfunkhouser/httpdumper/cmd/httpdumper@latest
go: downloading github.com/cfunkhouser/httpdumper v0.1.1
$ which httpdumper
/home/whatever/bin/httpdumper
$ httpdumper --help
 -address string
    	Bind address in host:port format. (default ":8080")
$ httpdumper --address 127.0.0.1:9000
INFO[0000] Server starting                               address="127.0.0.1:9000"
```

Or, you can run it using Docker:

```console
$ docker run --rm -p 8080:8080 ghcr.io/cfunkhouser/httpdumper/httpdumper:latest
Unable to find image 'ghcr.io/cfunkhouser/httpdumper/httpdumper:latest' locally
latest: Pulling from cfunkhouser/httpdumper/httpdumper
bd8f6a7501cc: Pull complete
44718e6d535d: Pull complete
efe9738af0cb: Pull complete
f37aabde37b8: Pull complete
c4c446e03742: Pull complete
03fadb054608: Pull complete
46075c0b6765: Pull complete
b8344ac560e4: Pull complete
561d809e37ca: Pull complete
Digest: sha256:c4ab8c77e79e858740a3db44785435165fed48a1de2f24585110a6471944d2e4
Status: Downloaded newer image for ghcr.io/cfunkhouser/httpdumper/httpdumper:latest
time="2021-04-21T19:09:21Z" level=info msg="Server starting" address=":8080"
```

## Client-side Logging

The very basic usage follows. You will likely want to tune your logger.

```go
c := &http.Client{
    Transport: httpdumper.DefaultTransport(),
}
```
