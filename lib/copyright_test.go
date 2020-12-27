package lib

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testConfigCopyright struct {
	showCopyrightText bool
}

func (c *testConfigCopyright) ShowCopyrightText() bool {
	return c.showCopyrightText
}

func TestGetCopyrightText(t *testing.T) {
	cfg := &testConfigCopyright{
		showCopyrightText: false,
	}

	copyright := GetCopyrightText(cfg)

	currentTime := time.Now()

	expected := `
	MIT License
	
	Copyright (c) %d Anton Sannikov`

	assert.Equal(t, fmt.Sprintf(expected, currentTime.Year()), copyright)

	cfg.showCopyrightText = true

	copyright = GetCopyrightText(cfg)

	assert.NotEqual(t, fmt.Sprintf(expected, currentTime.Year()), copyright)
}
