package proto

const (
	LowLevelProtocol = "tcp"

	SRS_A = 'A'
	SRS_B = 'B'
	SRS_C = 'C'

	RS_LengthA = 1
	RS_LengthB = 2
	RS_LengthC = 4

	RequestMessageMinLength = 3
	RequestMessageLengthA   = 255
	RequestMessageLengthB   = 65535

	ResponseMessageMinLength = 3
	ResponseMessageLengthA   = 255
	ResponseMessageLengthB   = 65535
	ResponseMessageLengthC   = 4294967295
)
