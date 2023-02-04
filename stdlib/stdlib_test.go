package stdlib_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"test/stdlib"
)

type badBuffer struct {
	buf bytes.Buffer
	err error
}

func (b *badBuffer) Write(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}

	return b.buf.Write(p)
}

func (b *badBuffer) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

var (
	errMock         = errors.New("error")
	errMockCopyPath = &fs.PathError{
		Path: "/dev/pmtx",
		Err:  errMock,
	}
)

func TestCeack(t *testing.T) {
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

			err := stdlib.Run(buf, tt.args.command, tt.args.args)
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
