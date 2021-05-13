package config

var InsertFailed = "INSERT_FAILED"
var InsertSuccess = "INSERT_SUCCESS"
var UpdateSuccess = "UPDATE_SUCCESS"
var DeleteSuccess = "DELETE_SUCCESS"

var Dropshipper = "HANDLER_DROPSHIP_VIS/MIN"

var ValidServiceFeeTypes = []string{"SERVICE_FEE_PREMIUM", "SERVICE_FEE"}
var UserTypeSeller = "Seller"
var UserTypeDropshipper = "Dropshipper"
var UserTypeAdmin = "Admin"

// Delivery Statuses
var DeliveryStatusPendingApproval = "Pending Approval"
var DeliveryStatusProposed = "Proposed"
var DeliveryStatusAccepted = "Accepted"
var DeliveryStatusFulfilled = "Fulfilled"
var DeliveryStatusDelivered = "Delivered"
var DeliveryStatusReturned = "Returned"
var DeliveryStatusRejected = "Rejected"
var DeliveryStatusVoided = "Voided"