// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"errors"
)

// ErrNotInitialized is returned by SQLConnectionProducer.Connection if
// SQLConnectionProducer.Init has not yet been called. It is also used in an
// analogous way in the forked copies of SQLConnectionProducer used in some
// database plugins.
var ErrNotInitialized = errors.New("connection has not been initialized")
