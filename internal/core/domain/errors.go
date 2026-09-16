package domain

import "errors"

var ErrStoreAlreadyInitialized = errors.New("store already has an owner")
var ErrInsufficientStock = errors.New("stock quantity cannot be negative")
var ErrCategoryNotInStore = errors.New("category does not belong to store")
var ErrInvalidProduct = errors.New("invalid product")
var ErrInvalidCategory = errors.New("invalid category")
var ErrInvalidInventoryMovement = errors.New("invalid inventory movement")
