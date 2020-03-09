package ft232h

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	BlessOpenFlag()

	//flag.Parse()
	os.Exit(m.Run())
}

func TestNewFT232H(t *testing.T) {

	ft, err := NewFT232H()
	if nil != err {
		t.Fatalf("could not open device: %v", err)
	}

	if nil != ft.open {
		// exercise each of the open masks individually
		arg := map[string]string{}
		ft.open.flag.Visit(func(f *flag.Flag) {
			arg[f.Name] = f.Value.String()
		})
		for f, v := range arg {
			if err = ft.Close(); nil != err {
				t.Fatalf("could not close device: %v", err)
			}
			ft, err = NewFT232HWithFlag([]string{"-" + f, v}, false)
			if nil != err {
				t.Fatalf("could not open device: %v", err)
			}
		}
	}

	err = ft.Close()
	if nil != err {
		t.Fatalf("could not close device: %v", err)
	}

}
