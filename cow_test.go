package cowsay

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCow_Clone(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		from *Cow
		want *Cow
	}{
		{
			name: "without options",
			opts: []Option{},
			from: func() *Cow {
				cow, _ := New()
				return cow
			}(),
			want: func() *Cow {
				cow, _ := New()
				return cow
			}(),
		},
		{
			name: "with some options",
			opts: []Option{},
			from: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BalloonWidth(60),
				)
				return cow
			}(),
			want: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BalloonWidth(60),
				)
				return cow
			}(),
		},
		{
			name: "clone and some options",
			opts: []Option{
				Thinking(),
				Thoughts('o'),
			},
			from: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BalloonWidth(60),
				)
				return cow
			}(),
			want: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BalloonWidth(60),
					Thinking(),
					Thoughts('o'),
				)
				return cow
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.want.Clone(tt.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(Cow{}),
				cmpopts.IgnoreFields(Cow{}, "buf")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}

	t.Run("random", func(t *testing.T) {
		cow, _ := New(
			Type(""),
			Thinking(),
			Thoughts('o'),
			Eyes("xx"),
			Tongue("u"),
			Random(),
		)

		cloned, _ := cow.Clone()

		if diff := cmp.Diff(cow, cloned,
			cmp.AllowUnexported(Cow{}),
			cmpopts.IgnoreFields(Cow{}, "buf")); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("error", func(t *testing.T) {
		cow, err := New()
		if err != nil {
			t.Fatal(err)
		}

		wantErr := errors.New("error")
		_, err = cow.Clone(func(*Cow) error {
			return wantErr
		})
		if wantErr != err {
			t.Fatalf("want %v, but got %v", wantErr, err)
		}
	})
}

func Test_adjustTo2Chars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "  ",
		},
		{
			name: "1 character",
			s:    "1",
			want: "1 ",
		},
		{
			name: "2 characters",
			s:    "12",
			want: "12",
		},
		{
			name: "3 characters",
			s:    "123",
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := adjustTo2Chars(tt.s); got != tt.want {
				t.Errorf("adjustTo2Chars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFound_Error(t *testing.T) {
	file := "test"
	n := &NotFound{
		Cowfile: file,
	}
	want := fmt.Sprintf("not found %q cowfile", file)
	if want != n.Error() {
		t.Fatalf("want %q but got %q", want, n.Error())
	}
}
