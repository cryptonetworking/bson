package bson

type SubType byte

const (
	SubtypeGenericBinary        SubType = 0x00
	SubtypeFunction             SubType = 0x01
	SubtypeBinary               SubType = 0x02
	SubtypeUUIDold              SubType = 0x03
	SubtypeUUID                 SubType = 0x04
	SubtypeMD5                  SubType = 0x05
	SubtypeEncryptedBSONValue   SubType = 0x06
	SubtypeCompressedBSONColumn SubType = 0x07
	SubtypeUserDefined          SubType = 0x80
)
