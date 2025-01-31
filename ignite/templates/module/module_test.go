package module

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtoPackageName(t *testing.T) {
	cases := []struct {
		name   string
		app    string
		module string
		want   string
	}{
		{
			name:   "name",
			app:    "ignite",
			module: "test",
			want:   "ignite.test",
		},
		{
			name:   "path",
			app:    "ignite/cli",
			module: "test",
			want:   "ignite.cli.test",
		},
		{
			name:   "path with dash",
			app:    "ignite-hq/cli",
			module: "test",
			want:   "ignitehq.cli.test",
		},
		{
			name:   "path with number prefix",
			app:    "0ignite/cli",
			module: "test",
			want:   "_0ignite.cli.test",
		},
		{
			name:   "path with number prefix and dash",
			app:    "0ignite-hq/cli",
			module: "test",
			want:   "_0ignitehq.cli.test",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ProtoPackageName(tt.app, tt.module))
		})
	}
}
