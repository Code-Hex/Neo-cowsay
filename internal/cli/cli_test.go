package cli

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	clis := []struct {
		name     string
		thinking bool
	}{
		{
			name:     "cowsay",
			thinking: false,
		},
		{
			name:     "cowthink",
			thinking: true,
		},
	}
	for _, cli := range clis {
		cli := cli
		t.Run(cli.name, func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name     string
				phrase   string
				argv     []string
				testfile string
			}{
				{
					name:     "-n",
					phrase:   "foo\nbar\nbaz",
					argv:     []string{"-n"},
					testfile: "n_option.txt",
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					var stdout bytes.Buffer
					c := &CLI{
						Thinking: cli.thinking,
						writer:   &stdout,
						reader:   strings.NewReader(tt.phrase),
					}
					exit := c.Run(tt.argv)
					if exit != 0 {
						t.Fatalf("unexpected exit code: %d", exit)
					}
					testpath := filepath.Join("testdata", testDir(cli.thinking), tt.testfile)
					content, err := ioutil.ReadFile(testpath)
					if err != nil {
						t.Fatal(err)
					}
					want := string(content)
					if got := stdout.String(); want != got {
						t.Errorf("want\n%s\n-----got\n%s\n", want, got)
					}
				})
			}
		})
	}
}

func testDir(thinking bool) string {
	if thinking {
		return "cowthink"
	}
	return "cowsay"
}
