package db

import (
	"testing"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := [...]struct {
		name     string
		host     string
		password string
		migrate  string
		error    string
	}{
		{
			name:     "success without migration",
			host:     "localhost",
			password: "test-password",
			migrate:  "false",
		},
		{
			name:     "success with migration",
			host:     "localhost",
			password: "test-password",
			migrate:  "true",
		},
		{
			name:     "wrong host",
			host:     "z",
			password: "test-password",
			migrate:  "false",
			error:    "no such host",
		},
		{
			name:     "wrong password",
			host:     "localhost",
			password: "wrong-password",
			migrate:  "false",
			error:    "pq: password authentication failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := config.NewMock(1, 1, "", tc.password, tc.host, tc.migrate)
			assert.NoError(t, err)

			err = New()
			if err != nil {
				assert.NotEmpty(t, tc.error)
				assert.ErrorContains(t, err, tc.error)
				return
			}

			assert.Empty(t, tc.error)
		})
	}
}
