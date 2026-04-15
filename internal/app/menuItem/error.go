package menuItem

import "errors"

var errMenuCategoryNotFound = errors.New("menu-category-not-found")
var errPrinterNotFound = errors.New("printer-not-found")
var errMenuItemNotFound = errors.New("menu-item-not-found")
var errMenuItemSameTitleAlreadyExists = errors.New("menu-item-same-title-already-exists")
