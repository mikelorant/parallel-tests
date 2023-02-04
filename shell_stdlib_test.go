package shell_test

import (
	"io"
	"strings"
	"testing"

	"test"
	"github.com/stretchr/testify/assert"
)

func TestStdlib(t *testing.T) {
	t.Parallel()

	type args struct {
		command string
		args    []string
		err     error
	}
	type want struct {
		err string
		out string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid",
			args: args{
				command: "echo",
				args:    []string{"test"},
			},
			want: want{
				out: "test",
			},
		},
		{
			name: "invalid",
			args: args{
				command: "invalid",
			},
			want: want{
				err: "executable file not found in $PATH",
			},
		},
		{
			name: "empty",
			want: want{
				err: "no command",
			},
		},
		{
			name: "error_copy",
			args: args{
				command: "echo",
				args:    []string{"test"},
				err:     errMock,
			},
			want: want{
				err: "unable to copy output: error",
			},
		},
		{
			name: "error_copy_path",
			args: args{
				command: "echo",
				args:    []string{"test"},
				err:     errMockCopyPath,
			},
			want: want{
				err: "unable to copy output: /dev/pmtx: error",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &badBuffer{
				err: tt.args.err,
			}

			err := shell.Stdlib(buf, tt.args.command, tt.args.args)
			if tt.want.err != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.want.err)
				return
			}
			assert.NoError(t, err)

			out, _ := io.ReadAll(buf)
			assert.Equal(t, tt.want.out, strings.TrimSpace(string(out)))
		})
	}
}
