# iot-server

Server for [sula7/iot-reader](https://github.com/sula7/iot-reader). Listens TCP host:port, accepts packets.

### Env vars

````shell
LISTEN_ADDRESS=host:port # address listen to
LOG_LEVEL=debug # default is info
````

### Packet structure

#### Header

| Protocol version | Packet type | Reserved | Reserved |
|:----------------:|-------------|----------|----------|
|        1         | 10/11/20    | 0        | 0        |

Protocol versions:

* 1 - current version

In case of receiving unknown protocol version **TBD** (currently accepts any)

Packet types:

* 10 - ping (incoming) responses immediately with pong
* 11 - pong (response)
* 20 - data transfer (bi-directional)

In case of receiving unknown packet type the server will log an error and no response expected.

#### Body

Consists of 10 bytes. Structure **TBD** (currently not implemented)
