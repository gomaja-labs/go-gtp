package ie_test

import (
	"errors"
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func ExampleIE_Err() {
	invalid := ie.NewIMSI("1234567890123456")
	fmt.Println(errors.Is(invalid.Err(), ie.ErrInvalidLength))

	var missing *ie.IE
	_, err := missing.IMSI()
	fmt.Println(errors.Is(err, ie.ErrIENotFound))

	valid := ie.NewIMSI("001010123456789")
	if err := valid.Err(); err != nil {
		fmt.Println(err)
		return
	}
	imsi, err := valid.IMSI()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(imsi)
	// Output:
	// true
	// true
	// 001010123456789
}
