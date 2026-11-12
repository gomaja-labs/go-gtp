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
		size, ok := tvLengthMap[int(i.Type)]
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
			// TS 29.060 §7.7.27: PPP (1) and Non-IP (2) have no address.
			if (p[1] != 1 && p[1] != 2) || n != 2 {
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
			// TS 29.060 §7.7.27: either family may be static independently.
			if n != 2 && n != 6 && n != 18 && n != 22 {
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

	case QoSProfile:
		// TS 29.060 §7.7.34: allocation/retention priority and the base QoS fields.
		if n < 4 {
			return ErrInvalidLength
		}
	case IMEISV:
		// TS 29.060 §7.7.53: eight octets containing 15 or 16 digits.
		if n != 8 {
			return ErrInvalidLength
		}
		s, err := utils.DecodeBCDDigits(p)
		if err != nil {
			return err
		}
		if len(s) != 15 && len(s) != 16 {
			return ErrInvalidLength
		}
	case UserLocationInformation:
		// TS 29.060 §7.7.51: each supported geographical location occupies seven octets.
		if n < 8 {
			return ErrInvalidLength
		}
		if p[0] > 2 {
			return ErrMalformed
		}
		_, _, err := utils.DecodePLMN(p[1:4])
		return err
	case MSTimeZone:
		// TS 29.060 §7.7.52: two octets; reserved daylight value is not a valid value.
		if n != 2 {
			return ErrInvalidLength
		}
		if p[0]>>4 > 9 || p[1]&3 == 3 {
			return ErrMalformed
		}
	case ULITimestamp:
		if n != 4 {
			return ErrInvalidLength
		}
	case AuthenticationQuintuplet:
		// TS 29.060 §7.7.35: variable XRES and AUTN are each prefixed by a length.
		if n < 50 {
			return ErrInvalidLength
		}
		off := 49 + int(p[16])
		if off >= n || int(p[off]) != n-off-1 {
			return ErrInvalidLength
		}
	case ExtensionHeaderTypeList:
		// TS 29.060 §7.7.40: the count is one octet.
		if n > 255 {
			return ErrInvalidLength
		}
	case ProtocolConfigurationOptions:
		_, err := ParsePCOPayload(p)
		return err

	}
	return nil
}
