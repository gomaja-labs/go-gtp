package ie

import "encoding/binary"

// NewWithTypeExtension creates a type-254 IE with the extension preceding content.
func NewWithTypeExtension(ext uint16, instance uint8, content []byte) *IE {
	// TS 29.274 §8.2.1A: extensions are binary values from 256 through 65535.
	if ext < 256 {
		return newErrorIE(SpecialIETypeForIETypeExtension, ErrMalformed)
	}
	if len(content) > 0xffff-2 {
		return newErrorIE(SpecialIETypeForIETypeExtension, ErrInvalidLength)
	}
	payload := make([]byte, 2+len(content))
	binary.BigEndian.PutUint16(payload, ext)
	copy(payload[2:], content)
	return New(SpecialIETypeForIETypeExtension, instance, payload)
}

// TypeExtension returns the extended type, without interpreting its content.
func (i *IE) TypeExtension() (uint16, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != SpecialIETypeForIETypeExtension {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return 0, ErrInvalidLength
	}
	ext := binary.BigEndian.Uint16(i.Payload[:2])
	// TS 29.274 §8.2.1A: values below 256 are not extended types.
	if ext < 256 {
		return 0, ErrMalformed
	}
	return ext, nil
}
