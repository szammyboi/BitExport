# BitExport
Exporting Bytes in Golang with Bit-Packing

Current interface:
```go
type Example struct {
	a bool   `bits:"1"`
	b bool   `bits:"1"`
	c bool   `bits:"1"`
	d bool   `bits:"1"`
	e uint16 `bits:"12"`
	f uint8  `bits:"4"`
}

test := Example {
	a: 0,
	b: 1,
	c: 1,
	d: 0,
	e: 255,
}

// Encode the struct into 2 packed bytes
bytes := BitExport.ToBytes(test)

// Decode the bytes back into a struct
var decoded Example
BitExport.FromBytes(bytes, &decoded)

```
