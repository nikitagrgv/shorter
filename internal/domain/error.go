package domain

import "errors"

var ErrIdCollision = errors.New("link short id collision")
var ErrShortCodeCollision = errors.New("link short code collision")
var ErrTokenCollision = errors.New("link token collision")
var ErrLinkNotFound = errors.New("link not found")
