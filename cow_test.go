package cowsay

import (
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
					BallonWidth(60),
				)
				return cow
			}(),
			want: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BallonWidth(60),
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
					BallonWidth(60),
				)
				return cow
			}(),
			want: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BallonWidth(60),
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
}
