// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

const (
	cgilen    int = 7
	sailen    int = 7
	railen    int = 7
	tailen    int = 5
	ecgilen   int = 7
	lailen    int = 5
	menbilen  int = 6
	emenbilen int = 6
)

const (
	uliFlagCGI    uint8 = 0x01
	uliFlagSAI    uint8 = 0x02
	uliFlagRAI    uint8 = 0x04
	uliFlagTAI    uint8 = 0x08
	uliFlagECGI   uint8 = 0x10
	uliFlagLAI    uint8 = 0x20
	uliFlagMENBI  uint8 = 0x40
	uliFlagEMENBI uint8 = 0x80
)

type PLMN struct {
	MCC string
	MNC string
}

// CGI represents a CGI, which is defined to be used as a field of UserLocationInformation IE.
type CGI struct {
	*PLMN
	LAC uint16
	CI  uint16
}

// NewCGI creates a new CGI.
func NewCGI(mcc, mnc string, lac, ci uint16) *CGI {
	return &CGI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		LAC:  lac,
		CI:   ci,
	}
}

// HasCGI checks if UserLocationInformationFields has CGI.
func (f *UserLocationInformationFields) HasCGI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagCGI != 0
}

// SAI represents a SAI, which is defined to be used as a field of UserLocationInformation IE.
type SAI struct {
	*PLMN
	LAC uint16
	SAC uint16
}

// NewSAI creates a new SAI.
func NewSAI(mcc, mnc string, lac, sac uint16) *SAI {
	return &SAI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		LAC:  lac,
		SAC:  sac,
	}
}

// HasSAI checks if UserLocationInformationFields has SAI.
func (f *UserLocationInformationFields) HasSAI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagSAI != 0
}

// RAI represents a RAI, which is defined to be used as a field of UserLocationInformation IE.
type RAI struct {
	*PLMN
	LAC uint16
	RAC uint8
}

// NewRAI creates a new RAI.
func NewRAI(mcc, mnc string, lac uint16, rac uint8) *RAI {
	return &RAI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		LAC:  lac,
		RAC:  rac,
	}
}

// HasRAI checks if UserLocationInformationFields has RAI.
func (f *UserLocationInformationFields) HasRAI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagRAI != 0
}

// TAI represents a TAI, which is defined to be used as a field of UserLocationInformation IE.
type TAI struct {
	*PLMN
	TAC uint16
}

// NewTAI creates a new TAI.
func NewTAI(mcc, mnc string, tac uint16) *TAI {
	return &TAI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		TAC:  tac,
	}
}

// HasTAI checks if UserLocationInformationFields has TAI.
func (f *UserLocationInformationFields) HasTAI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagTAI != 0
}

// ECGI represents a ECGI, which is defined to be used as a field of UserLocationInformation IE.
type ECGI struct {
	*PLMN
	ECI uint32
}

// NewECGI creates a new ECGI.
func NewECGI(mcc, mnc string, eci uint32) *ECGI {
	return &ECGI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		ECI:  eci,
	}
}

// HasECGI checks if UserLocationInformationFields has ECGI.
func (f *UserLocationInformationFields) HasECGI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagECGI != 0
}

// LAI represents a LAI, which is defined to be used as a field of UserLocationInformation IE.
type LAI struct {
	*PLMN
	LAC uint16
}

// NewLAI creates a new LAI.
func NewLAI(mcc, mnc string, lac uint16) *LAI {
	return &LAI{
		PLMN: &PLMN{MCC: mcc, MNC: mnc},
		LAC:  lac,
	}
}

// HasLAI checks if UserLocationInformationFields has LAI.
func (f *UserLocationInformationFields) HasLAI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagLAI != 0
}

// MENBI represents a MENBI, which is defined to be used as a field of UserLocationInformation IE.
type MENBI struct {
	*PLMN
	MENBI uint32
}

// NewMENBI creates a new MENBI.
func NewMENBI(mcc, mnc string, menbi uint32) *MENBI {
	return &MENBI{
		PLMN:  &PLMN{MCC: mcc, MNC: mnc},
		MENBI: menbi,
	}
}

// HasMENBI checks if UserLocationInformationFields has MENBI.
func (f *UserLocationInformationFields) HasMENBI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagMENBI != 0
}

// EMENBI represents a EMENBI, which is defined to be used as a field of UserLocationInformation IE.
type EMENBI struct {
	*PLMN
	EMENBI uint32
	Short  bool
}

// NewEMENBI creates a new EMENBI.
func NewEMENBI(mcc, mnc string, menbi uint32) *EMENBI {
	return &EMENBI{
		PLMN:   &PLMN{MCC: mcc, MNC: mnc},
		EMENBI: menbi,
	}
}

// HasEMENBI checks if UserLocationInformationFields has EMENBI.
func (f *UserLocationInformationFields) HasEMENBI() bool {
	if f == nil {
		return false
	}
	return f.Flags&uliFlagEMENBI != 0
}

// NewUserLocationInformation creates a new UserLocationInformation IE from
// the structs defined in gtpv2/ie package. Give nil for unnecessary values.
func NewUserLocationInformation(cgi *CGI, sai *SAI, rai *RAI, tai *TAI, ecgi *ECGI, lai *LAI, menbi *MENBI, emenbi *EMENBI) *IE {
	v := NewUserLocationInformationFields(cgi, sai, rai, tai, ecgi, lai, menbi, emenbi)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(UserLocationInformation, err)
	}

	return New(UserLocationInformation, 0x00, b)
}

// UserLocationInformation returns UserLocationInformation in UserLocationInformationFields type
// if the type of IE matches.
func (i *IE) UserLocationInformation() (*UserLocationInformationFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case UserLocationInformation:
		return ParseUserLocationInformationFields(i.Payload)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// UserLocationInformationFields is a set of fields in UserLocationInformation IE.
//
// The contained fields are of type struct, as they are too complex to handle with
// existing (standard) types in Go.
type UserLocationInformationFields struct {
	Flags uint8
	*CGI
	*SAI
	*RAI
	*TAI
	*ECGI
	*LAI
	*MENBI
	*EMENBI
}

// NewUserLocationInformationFields creates a new NewUserLocationInformationFields.
func NewUserLocationInformationFields(cgi *CGI, sai *SAI, rai *RAI, tai *TAI, ecgi *ECGI, lai *LAI, menbi *MENBI, emenbi *EMENBI) *UserLocationInformationFields {
	f := &UserLocationInformationFields{}

	if cgi != nil {
		f.Flags |= uliFlagCGI
		f.CGI = cgi
	}
	if sai != nil {
		f.Flags |= uliFlagSAI
		f.SAI = sai
	}
	if rai != nil {
		f.Flags |= uliFlagRAI
		f.RAI = rai
	}
	if tai != nil {
		f.Flags |= uliFlagTAI
		f.TAI = tai
	}
	if ecgi != nil {
		f.Flags |= uliFlagECGI
		f.ECGI = ecgi
	}
	if lai != nil {
		f.Flags |= uliFlagLAI
		f.LAI = lai
	}
	if menbi != nil {
		f.Flags |= uliFlagMENBI
		f.MENBI = menbi
	}
	if emenbi != nil {
		f.Flags |= uliFlagEMENBI
		f.EMENBI = emenbi
	}

	return f
}

// Marshal serializes UserLocationInformationFields.
func (f *UserLocationInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes UserLocationInformationFields.
func (f *UserLocationInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	// TS 29.274 §§8.21.7-8.21.8: at most one eNodeB identity.
	if f.Flags&0xc0 == 0xc0 {
		return ErrMalformed
	}
	if f.EMENBI != nil && f.EMENBI.Short && f.EMENBI.EMENBI > 0x3ffff {
		return ErrMalformed
	}
	if f.ECGI != nil && f.ECGI.ECI > 0xfffffff || f.MENBI != nil && f.MENBI.MENBI > 0xfffff || f.EMENBI != nil && f.EMENBI.EMENBI > 0x1fffff {
		return ErrMalformed
	}
	if f.Flags&uliFlagCGI != 0 && (f.CGI == nil || f.CGI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagSAI != 0 && (f.SAI == nil || f.SAI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagRAI != 0 && (f.RAI == nil || f.RAI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagTAI != 0 && (f.TAI == nil || f.TAI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagECGI != 0 && (f.ECGI == nil || f.ECGI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagLAI != 0 && (f.LAI == nil || f.LAI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagMENBI != 0 && (f.MENBI == nil || f.MENBI.PLMN == nil) {
		return ErrMalformed
	}
	if f.Flags&uliFlagEMENBI != 0 && (f.EMENBI == nil || f.EMENBI.PLMN == nil) {
		return ErrMalformed
	}

	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags
	offset := 1

	if has1stBit(f.Flags) { // CGI
		if l < offset+cgilen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.CGI.MCC, f.CGI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint16(b[offset+3:offset+5], f.CGI.LAC)
		binary.BigEndian.PutUint16(b[offset+5:offset+7], f.CGI.CI)
		offset += cgilen
	}
	if has2ndBit(f.Flags) { // SAI
		if l < offset+sailen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.SAI.MCC, f.SAI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint16(b[offset+3:offset+5], f.SAI.LAC)
		binary.BigEndian.PutUint16(b[offset+5:offset+7], f.SAI.SAC)
		offset += sailen
	}
	if has3rdBit(f.Flags) { // RAI
		if l < offset+railen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.RAI.MCC, f.RAI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint16(b[offset+3:offset+5], f.RAI.LAC)
		// TS 29.274 §8.21.3: RAC followed by 0xff.
		b[offset+5], b[offset+6] = f.RAI.RAC, 0xff
		offset += railen
	}
	if has4thBit(f.Flags) { // TAI
		if l < offset+tailen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.TAI.MCC, f.TAI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint16(b[offset+3:offset+5], f.TAI.TAC)
		offset += tailen
	}
	if has5thBit(f.Flags) { // ECGI
		if l < offset+ecgilen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.ECGI.MCC, f.ECGI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint32(b[offset+3:offset+7], f.ECGI.ECI)
		offset += ecgilen
	}
	if has6thBit(f.Flags) { // LAI
		if l < offset+lailen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.LAI.MCC, f.LAI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		binary.BigEndian.PutUint16(b[offset+3:offset+5], f.LAI.LAC)
		offset += lailen
	}
	if has7thBit(f.Flags) { // MENBI
		if l < offset+menbilen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.MENBI.MCC, f.MENBI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		encoded, err := utils.Uint32To24(f.MENBI.MENBI)
		if err != nil {
			return err
		}
		copy(b[offset+3:offset+6], encoded)
		offset += menbilen
	}
	if has8thBit(f.Flags) { // EMENBI
		if l < offset+emenbilen {
			return io.ErrUnexpectedEOF
		}

		plmn, err := utils.EncodePLMN(f.EMENBI.MCC, f.EMENBI.MNC)
		if err != nil {
			return err
		}
		copy(b[offset:offset+3], plmn)
		encoded, err := utils.Uint32To24(f.EMENBI.EMENBI)
		if err != nil {
			return err
		}
		copy(b[offset+3:offset+6], encoded)
		if f.EMENBI.Short {
			b[offset+3] |= 0x80
		}
	}

	return nil
}

// ParseUserLocationInformationFields decodes UserLocationInformationFields.
func ParseUserLocationInformationFields(b []byte) (*UserLocationInformationFields, error) {
	f := &UserLocationInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into UserLocationInformationFields.
func (f *UserLocationInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = UserLocationInformationFields{}
	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0]
	if f.Flags&0xc0 == 0xc0 {
		return ErrMalformed
	}
	offset := 1

	var err error
	if has1stBit(f.Flags) { // CGI
		if l < offset+cgilen {
			return io.ErrUnexpectedEOF
		}

		f.CGI = &CGI{PLMN: &PLMN{}}
		f.CGI.PLMN.MCC, f.CGI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.CGI.LAC = binary.BigEndian.Uint16(b[offset+3 : offset+5])
		f.CGI.CI = binary.BigEndian.Uint16(b[offset+5 : offset+7])
		offset += cgilen
	}
	if has2ndBit(f.Flags) { // SAI
		if l < offset+sailen {
			return io.ErrUnexpectedEOF
		}

		f.SAI = &SAI{PLMN: &PLMN{}}
		f.SAI.PLMN.MCC, f.SAI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.SAI.LAC = binary.BigEndian.Uint16(b[offset+3 : offset+5])
		f.SAI.SAC = binary.BigEndian.Uint16(b[offset+5 : offset+7])
		offset += sailen
	}
	if has3rdBit(f.Flags) { // RAI
		if l < offset+railen {
			return io.ErrUnexpectedEOF
		}

		f.RAI = &RAI{PLMN: &PLMN{}}
		f.RAI.PLMN.MCC, f.RAI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.RAI.LAC = binary.BigEndian.Uint16(b[offset+3 : offset+5])
		f.RAI.RAC = b[offset+5]
		offset += railen
	}
	if has4thBit(f.Flags) { // TAI
		if l < offset+tailen {
			return io.ErrUnexpectedEOF
		}

		f.TAI = &TAI{PLMN: &PLMN{}}
		f.TAI.PLMN.MCC, f.TAI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.TAI.TAC = binary.BigEndian.Uint16(b[offset+3 : offset+5])
		offset += tailen
	}
	if has5thBit(f.Flags) { // ECGI
		if l < offset+ecgilen {
			return io.ErrUnexpectedEOF
		}

		f.ECGI = &ECGI{PLMN: &PLMN{}}
		f.ECGI.PLMN.MCC, f.ECGI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.ECGI.ECI = binary.BigEndian.Uint32(b[offset+3:offset+7]) & 0xfffffff
		offset += ecgilen
	}
	if has6thBit(f.Flags) { // LAI
		if l < offset+lailen {
			return io.ErrUnexpectedEOF
		}

		f.LAI = &LAI{PLMN: &PLMN{}}
		f.LAI.PLMN.MCC, f.LAI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.LAI.LAC = binary.BigEndian.Uint16(b[offset+3 : offset+5])
		offset += lailen
	}
	if has7thBit(f.Flags) { // MENBI
		if l < offset+menbilen {
			return io.ErrUnexpectedEOF
		}

		f.MENBI = &MENBI{PLMN: &PLMN{}}
		f.MENBI.PLMN.MCC, f.MENBI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.MENBI.MENBI = utils.Uint24To32(b[offset+3:offset+6]) & 0xfffff
		offset += menbilen
	}
	if has8thBit(f.Flags) { // EMENBI
		if l < offset+emenbilen {
			return io.ErrUnexpectedEOF
		}

		f.EMENBI = &EMENBI{PLMN: &PLMN{}}
		f.EMENBI.PLMN.MCC, f.EMENBI.PLMN.MNC, err = utils.DecodePLMN(b[offset : offset+3])
		if err != nil {
			return err
		}
		f.EMENBI.EMENBI = utils.Uint24To32(b[offset+3:offset+6]) & 0x1fffff
		f.EMENBI.Short = b[offset+3]&0x80 != 0
		if f.EMENBI.Short {
			f.EMENBI.EMENBI &= 0x3ffff
		}
	}

	return nil
}

// MarshalLen returns the serial length of UserLocationInformationFields in int.
func (f *UserLocationInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	l := 1

	if has1stBit(f.Flags) {
		l += cgilen
	}
	if has2ndBit(f.Flags) {
		l += sailen
	}
	if has3rdBit(f.Flags) {
		l += railen
	}
	if has4thBit(f.Flags) {
		l += tailen
	}
	if has5thBit(f.Flags) {
		l += ecgilen
	}
	if has6thBit(f.Flags) {
		l += lailen
	}
	if has7thBit(f.Flags) {
		l += menbilen
	}
	if has8thBit(f.Flags) {
		l += emenbilen
	}

	return l
}
