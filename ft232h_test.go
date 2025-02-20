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

	t.Run("Info", func(t *testing.T) {
		if ft.IsOpen() != ft.info.isOpen {
			t.Errorf("IsOpen() != info.isOpen")
		}
		t.Logf("FT232H IsOpen: %v", ft.IsOpen())
		if ft.Serial() != ft.info.serial {
			t.Errorf("Serial() != info.serial")
		}
		t.Logf("FT232H Serial: %s", ft.Serial())
		if ft.Desc() != ft.info.desc {
			t.Errorf("Description() != info.description")
		}
		t.Logf("FT232H Description: %s", ft.Desc())
		if ft.PID() != ft.info.pid {
			t.Errorf("PID() != info.pid")
		}
		t.Logf("FT232H PID: %d", ft.PID())
		if ft.VID() != ft.info.vid {
			t.Errorf("VID() != info.vid")
		}
		t.Logf("FT232H VID: %d", ft.VID())
		if ft.Index() != ft.info.index {
			t.Errorf("Index() != info.index")
		}
		t.Logf("FT232H Index: %d", ft.Index())
	})

	err = ft.Close()
	if nil != err {
		t.Fatalf("could not close device: %v", err)
	}

}
