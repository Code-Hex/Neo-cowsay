package cowsay

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCow_Clone(t *testing.T) {
	tests := []struct {
		name string
		want *Cow
	}{
		{
			name: "non option",
			want: func() *Cow {
				cow, _ := New()
				return cow
			}(),
		},
		{
			name: "with some options",
			want: func() *Cow {
				cow, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return cow
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.want.Clone()
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(Cow{}),
				cmpopts.IgnoreFields(Cow{}, "buf")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
