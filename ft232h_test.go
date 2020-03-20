package ft232h

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	BlessFlag()

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {

	if testing.Short() {
		t.Skipf("short: skipping FT232H open tests")
	}

	ft, err := New()
	if nil != err {
		t.Fatalf("could not open device: %v", err)
	}

	t.Logf("opened: %s", ft)

	if nil != ft.flag {
		// exercise each of the open masks individually
		arg := map[string]string{}
		ft.flag.Visit(func(f *flag.Flag) {
			arg[f.Name] = f.Value.String()
		})
		for f, v := range arg {
			if err = ft.Close(); nil != err {
				t.Fatalf("could not close device: %v", err)
			}
			ft, err = OpenFlag([]string{"-" + f, v}, false)
			if nil != err {
				t.Fatalf("could not open device: %v", err)
			}
			t.Logf("opened: %s", ft)
		}
	}

	err = ft.Close()
	if nil != err {
		t.Fatalf("could not close device: %v", err)
	}

}
