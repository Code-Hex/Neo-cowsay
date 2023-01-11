package cowsay

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCows(t *testing.T) {
	t.Run("no set COWPATH env", func(t *testing.T) {
		cowPaths, err := Cows()
		if err != nil {
			t.Fatal(err)
		}
		if len(cowPaths) != 1 {
			t.Fatalf("want 1, but got %d", len(cowPaths))
		}
		cowPath := cowPaths[0]
		if len(cowPath.CowFiles) == 0 {
			t.Fatalf("no cowfiles")
		}

		wantCowPath := &CowPath{
			Name:         "cows",
			LocationType: InBinary,
		}
		if diff := cmp.Diff(wantCowPath, cowPath,
			cmpopts.IgnoreFields(CowPath{}, "CowFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("set COWPATH env", func(t *testing.T) {
		cowpath := filepath.Join("testdata", "testdir")

		os.Setenv("COWPATH", cowpath)
		defer os.Unsetenv("COWPATH")

		cowPaths, err := Cows()
		if err != nil {
			t.Fatal(err)
		}
		if len(cowPaths) != 2 {
			t.Fatalf("want 2, but got %d", len(cowPaths))
		}

		wants := []*CowPath{
			{
				Name:         filepath.Join("testdata", "testdir"),
				LocationType: InDirectory,
			},
			{
				Name:         "cows",
				LocationType: InBinary,
			},
		}
		if diff := cmp.Diff(wants, cowPaths,
			cmpopts.IgnoreFields(CowPath{}, "CowFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}

		if len(cowPaths[0].CowFiles) != 1 {
			t.Fatalf("unexpected cowfiles len = %d, %+v",
				len(cowPaths[0].CowFiles), cowPaths[0].CowFiles,
			)
		}

		if cowPaths[0].CowFiles[0] != "test" {
			t.Fatalf("want %q but got %q", "test", cowPaths[0].CowFiles[0])
		}
	})

	t.Run("set COWPATH env", func(t *testing.T) {
		os.Setenv("COWPATH", "notfound")
		defer os.Unsetenv("COWPATH")

		_, err := Cows()
		if err == nil {
			t.Fatal("want error")
		}
	})

}

func TestCowPath_Lookup(t *testing.T) {
	t.Run("looked for cowfile", func(t *testing.T) {
		c := &CowPath{
			Name:         "basepath",
			CowFiles:     []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("test")
		if !ok {
			t.Errorf("want %v", ok)
		}
		want := &CowFile{
			Name:         "test",
			BasePath:     "basepath",
			LocationType: InBinary,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("no cowfile", func(t *testing.T) {
		c := &CowPath{
			Name:         "basepath",
			CowFiles:     []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("no cowfile")
		if ok {
			t.Errorf("want %v", !ok)
		}
		if got != nil {
			t.Error("want nil")
		}
	})
}

func TestCowFile_ReadAll(t *testing.T) {
	fromTestData := &CowFile{
		Name:         "test",
		BasePath:     filepath.Join("testdata", "testdir"),
		LocationType: InDirectory,
	}
	fromTestdataContent, err := fromTestData.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	fromBinary := &CowFile{
		Name:         "default",
		BasePath:     "cows",
		LocationType: InBinary,
	}
	fromBinaryContent, err := fromBinary.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fromTestdataContent, fromBinaryContent) {
		t.Fatalf("testdata\n%s\n\nbinary%s\n", string(fromTestdataContent), string(fromBinaryContent))
	}

}

const defaultSay = ` ________ 
< cowsay >
 -------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`

func TestSay(t *testing.T) {
	type args struct {
		phrase  string
		options []Option
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "default",
			args: args{
				phrase: "hello!",
			},
			wantFile: "default.cow",
			wantErr:  false,
		},
		{
			name: "nest",
			args: args{
				phrase: defaultSay,
				options: []Option{
					DisableWordWrap(),
				},
			},
			wantFile: "nest.cow",
			wantErr:  false,
		},
		{
			name: "error",
			args: args{
				phrase: "error",
				options: []Option{
					func(*Cow) error {
						return errors.New("error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Say(tt.args.phrase, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Say() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			filename := filepath.Join("testdata", tt.wantFile)
			content, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			got = strings.Replace(got, "\r", "", -1)               // for windows
			want := strings.Replace(string(content), "\r", "", -1) // for windows
			if want != got {
				t.Log(cmp.Diff([]byte(want), []byte(got)))
				t.Fatalf("want\n%s\n\ngot\n%s", want, got)
			}
		})
	}
}
