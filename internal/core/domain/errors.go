package domain

import "errors"

var ErrStoreAlreadyInitialized = errors.New("store already has an owner")
var ErrInsufficientStock = errors.New("stock quantity cannot be negative")
var ErrCategoryNotInStore = errors.New("category does not belong to store")
var ErrInvalidProduct = errors.New("invalid product")
var ErrInvalidCategory = errors.New("invalid category")
var ErrInvalidInventoryMovement = errors.New("invalid inventory movement")
var ErrInvalidStaff = errors.New("invalid staff member")
var ErrForbiddenRoleAssignment = errors.New("role assignment is not allowed")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrInvalidOrder = errors.New("invalid order")
var ErrProductNotFound = errors.New("product not found")
var ErrOrderNotVoidable = errors.New("order cannot be voided")
var ErrPaidOrderCannotVoid = errors.New("paid order requires a refund")
