package ie

import (
	"github.com/gomaja-labs/go-gtp/utils"
	"strings"
)

// validatePayload checks the value independently of a stored wire length.
func (i *IE) validatePayload() error {
	p := i.Payload
	n := len(p)
	// TS 29.060 Table 37; GSM 09.60 §7.9: TV lengths are fixed.
	if i.Type < 128 {
		size, ok := tvLengthMap[i.Type]
		if !ok {
			return &InvalidTypeError{Type: i.Type}
		}
		if n != size {
			return ErrInvalidLength
		}
	}
	if n > 65535 {
		return ErrInvalidLength
	}
	switch i.Type {
	case IMSI:
		// TS 29.060 §7.7.2; GSM 09.60 §7.9.2: terminal filler and at most 15 digits.
		s, err := utils.DecodeBCDDigits(p)
		if err != nil {
			return err
		}
		if len(s) == 0 || len(s) > 15 {
			return ErrInvalidLength
		}
	case RouteingAreaIdentity:
		_, _, err := utils.DecodePLMN(p)
		return err
	case MSISDN:
		if n < 2 {
			return ErrInvalidLength
		}
		s, err := utils.DecodeBCDDigits(p[1:])
		if err != nil {
			return err
		}
		if len(s) == 0 || len(s) > 15 {
			return ErrInvalidLength
		}
	case AccessPointName:
		// TS 23.003 §9.1: labels without a terminating zero, at most 100 octets.
		if n == 0 || n > 100 {
			return ErrInvalidLength
		}
		s, err := utils.DecodeFQDN(p)
		if err != nil {
			return err
		}
		if err := validateAPNName(s); err != nil {
			return err
		}
		encoded, err := utils.EncodeFQDN(s)
		if err != nil {
			return err
		}
		if len(encoded) != n || strings.HasSuffix(s, ".") {
			return ErrMalformed
		}
	case EndUserAddress:
		// TS 29.060 §7.7.27; GSM 09.60 §7.9.18: spare organization bits are ignored.
		if n < 2 {
			return ErrInvalidLength
		}
		if p[0]&15 == 0 {
			if p[1] != 1 || n != 2 {
				return ErrMalformed
			}
			return nil
		}
		if p[0]&15 != 1 {
			return ErrMalformed
		}
		switch p[1] {
		case 0x21:
			if n != 2 && n != 6 {
				return ErrInvalidLength
			}
		case 0x57:
			if n != 2 && n != 18 {
				return ErrInvalidLength
			}
		case 0x8d:
			if n != 2 && n != 22 {
				return ErrInvalidLength
			}
		default:
			return ErrMalformed
		}
	case GSNAddress, ChargingGatewayAddress:
		// TS 29.060 §7.7.32; GSM 09.60 §§7.9.23, 7.9.25: IPv4 or IPv6.
		if n != 4 && n != 16 {
			return ErrInvalidLength
		}
	case PrivateExtension:
		// TS 29.060 §7.7.46; GSM 09.60 §7.9.26: two-octet identifier precedes value.
		if n < 2 {
			return ErrInvalidLength
		}

	}
	return nil
}
