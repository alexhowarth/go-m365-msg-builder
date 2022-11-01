# M365 Message Builder

This is a Go implementation of the reverse engineered Xiaomi M365 Bluetooth Low Energy (BLE) Electric Scooter [protocol](https://github.com/CamiAlfa/M365-BLE-PROTOCOL).

## Usage

```go
// turn light on

m := m365message.NewMessage()
m.SetRW(m365message.WRITE)
m.SetPayload([]int{0x0002})
m.SetPosition(0x7D)
m.SetDirection(m365message.MASTER_TO_M365)
data, err := m.Build()
if err != nil {
    // handle error
}

fmt.Println(data)
// send data via BLE to Scooter!

```

## Protocol

A detailed explaination of the protocol by Camilo Ruiz can be found here:

https://github.com/CamiAlfa/M365-BLE-PROTOCOL/blob/master/protocolo

The protocol appears to have changed for scooters with firmware >= 1.5.x which requires an additional level of authentication with the scooter that this library does not provide.