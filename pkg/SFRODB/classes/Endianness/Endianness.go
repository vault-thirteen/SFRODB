package endianness

const (
	Endianness_BigEndian    = Endianness(1)
	Endianness_LittleEndian = Endianness(2)
)

const (
	ErrEndiannessIsUnknown = "endianness is unknown"
)

type Endianness byte
