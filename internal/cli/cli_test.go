package cli

import (
	"bytes"
	"fmt"
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
					name:     "ignore wordwrap option",
					phrase:   "foo\nbar\nbaz",
					argv:     []string{"-n"},
					testfile: "n_option.txt",
				},
				{
					name:     "tired option",
					phrase:   "tired",
					argv:     []string{"-t"},
					testfile: "t_option.txt",
				},
				{
					name:     "specifies width of the ballon is 3",
					phrase:   "foobarbaz",
					argv:     []string{"-W", "3"},
					testfile: "W_option.txt",
				},
				{
					name:     "borg mode",
					phrase:   "foobarbaz",
					argv:     []string{"-b"},
					testfile: "b_option.txt",
				},
				{
					name:     "dead mode",
					phrase:   "0xdeadbeef",
					argv:     []string{"-d"},
					testfile: "d_option.txt",
				},
				{
					name:     "greedy mode",
					phrase:   "give me money",
					argv:     []string{"-g"},
					testfile: "g_option.txt",
				},
				{
					name:     "paranoid mode",
					phrase:   "everyone hates me",
					argv:     []string{"-p"},
					testfile: "p_option.txt",
				},
				{
					name:     "stoned mode",
					phrase:   "I don't know",
					argv:     []string{"-s"},
					testfile: "s_option.txt",
				},
				{
					name:     "wired mode",
					phrase:   "Wanna Netflix and chill?",
					argv:     []string{"-w"},
					testfile: "wired_option.txt",
				},
				{
					name:     "youthful mode",
					phrase:   "I forgot my ID at home",
					argv:     []string{"-y"},
					testfile: "y_option.txt",
				},
				{
					name:     "eyes option",
					phrase:   "I'm not angry",
					argv:     []string{"-e", "^^"},
					testfile: "eyes_option.txt",
				},
				{
					name:     "tongue option",
					phrase:   "hungry",
					argv:     []string{"-T", ":"},
					testfile: "tongue_option.txt",
				},
				{
					name:     "-f tux",
					phrase:   "what is macOS?",
					argv:     []string{"-f", "tux"},
					testfile: "f_tux_option.txt",
				},
			}
			for _, tt := range tests {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					var stdout bytes.Buffer
					c := &CLI{
						Thinking: cli.thinking,
						stdout:   &stdout,
						stdin:    strings.NewReader(tt.phrase),
					}
					exit := c.Run(tt.argv)
					if exit != 0 {
						t.Fatalf("unexpected exit code: %d", exit)
					}
					testpath := filepath.Join("testdata", cli.name, tt.testfile)
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

			t.Run("program name", func(t *testing.T) {
				c := &CLI{Thinking: cli.thinking}
				if cli.name != c.program() {
					t.Fatalf("want %q, but got %q", cli.name, c.program())
				}
			})

			t.Run("not found cowfile", func(t *testing.T) {
				var stderr bytes.Buffer
				c := &CLI{
					Thinking: cli.thinking,
					stderr:   &stderr,
				}

				exit := c.Run([]string{"-f", "unknown"})
				if exit == 0 {
					t.Errorf("unexpected exit code: %d", exit)
				}
				want := fmt.Sprintf("%s: Could not find unknown cowfile!\n", cli.name)
				if want != stderr.String() {
					t.Errorf("want %q, but got %q", want, stderr.String())
				}
			})
		})
	}
}
