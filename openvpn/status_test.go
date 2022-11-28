package openvpn

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
	"time"
)

func TestParseFile(t *testing.T) {

	tests := []struct {
		name                string
		filename            string
		expectedUp          bool
		expectedClientCount uint
		expectedRouteCount  uint
		expectedError       error
	}{
		{
			name:                "single client",
			filename:            "test-samples/status-single.log",
			expectedUp:          true,
			expectedClientCount: 1,
			expectedRouteCount:  1,
			expectedError:       nil,
		}, {
			name:                "multiple client",
			filename:            "test-samples/status-multiple.log",
			expectedUp:          true,
			expectedClientCount: 13,
			expectedRouteCount:  13,
			expectedError:       nil,
		}, {
			name:                "no client",
			filename:            "test-samples/status-noclient.log",
			expectedUp:          true,
			expectedClientCount: 0,
			expectedRouteCount:  0,
			expectedError:       nil,
		}, {
			name:                "empty file",
			filename:            "test-samples/status-empty.log",
			expectedUp:          true,
			expectedClientCount: 0,
			expectedRouteCount:  0,
			expectedError:       &ParseError{},
		}, {
			name:                "invalid file",
			filename:            "test-samples/status-invalid.log",
			expectedUp:          true,
			expectedClientCount: 0,
			expectedRouteCount:  0,
			expectedError:       &ParseError{},
		}, {
			name:                "no status file",
			filename:            "test-samples/file-not-there.log",
			expectedUp:          false,
			expectedClientCount: 0,
			expectedRouteCount:  0,
			expectedError:       &fs.PathError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := time.Now()
			_, err := os.Stat(test.filename)
			if !os.IsNotExist(err) {
				err := os.Chtimes(test.filename, now, now)
				require.NoError(t, err)
			}

			status, err := ParseFile(test.filename)
			if test.expectedError == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expectedClientCount, status.ClientCount)
				assert.Equal(t, test.expectedRouteCount, status.RouteCount)
				assert.Equal(t, test.expectedUp, status.IsUp)
				assert.Equal(t, now.Unix(), status.LastModified.Unix())
			} else {
				assert.IsType(t, test.expectedError, err)
			}
		})
	}
}
