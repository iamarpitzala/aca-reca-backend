package util

// Error messages for HTTP responses
const (
	// Validation errors
	ErrInvalidClinicID      = "invalid clinic ID"
	ErrInvalidFormID        = "invalid form ID"
	ErrInvalidEntryID       = "invalid entry ID"
	ErrInvalidUserID        = "invalid user ID"
	ErrInvalidSessionID     = "invalid session ID"
	ErrInvalidAssociationID = "invalid association ID"
	ErrInvalidAccountTaxID  = "invalid account tax id"
	ErrInvalidAccountTypeID = "invalid account type id"
	ErrInvalidID            = "invalid id"
	ErrInvalidYearsBack     = "invalid yearsBack parameter"
	ErrInvalidYearsForward  = "invalid yearsForward parameter"
	ErrInvalidDateFormat    = "invalid date format. Use RFC3339 (2006-01-02T15:04:05Z07:00) or date (2006-01-02)"
	ErrDateRequired         = "date parameter is required"

	// Authentication errors
	ErrUnauthorized          = "unauthorized"
	ErrInvalidUserContext    = "invalid user context"
	ErrUserNotAuthenticated  = "user not authenticated"
	ErrRefreshTokenRequired  = "refreshToken is required"
	ErrProviderRequired      = "provider is required"
	ErrFailedToGenerateState = "failed to generate state"

	// Authorization errors
	ErrAccessDenied                    = "access denied: you do not have access to this clinic"
	ErrAccessDeniedOwnerRequired       = "access denied: owner access required"
	ErrAccessDeniedOwnerOnly           = "access denied: only the clinic owner can perform this action"
	ErrAccessDeniedOwnerCanAddUsers    = "access denied: only the clinic owner can add users"
	ErrAccessDeniedOwnerCanRemoveUsers = "access denied: only the clinic owner can remove users"
	ErrAccessDeniedOwnClinicsOnly      = "access denied: you can only view your own clinics"

	// Not found errors
	ErrClinicNotFound                = "clinic not found"
	ErrAOCNotFound                   = "coa not found"
	ErrUserClinicAssociationNotFound = "user-clinic association not found"
	ErrNetDetailsNotFound            = "net details not found for this entry"
	ErrGrossDetailsNotFound          = "gross details not found for this entry"
	ErrSessionNotFound               = "session not found or does not belong to user"
	ErrTransactionNotFound           = "transaction not found"

	// Request validation errors
	ErrIDsRequired       = "invalid request: ids required"
	ErrIDsMustNotBeEmpty = "ids must not be empty"

	// Upload errors
	ErrUploadServiceNotConfigured = "upload service is not configured"
	ErrMissingOrInvalidFile       = "missing or invalid file: "
	ErrImageTooLarge              = "image must be at most 10 MB"
	ErrDocumentTooLarge           = "document must be at most 20 MB"
	ErrInvalidImageType           = "invalid image type; allowed: JPEG, PNG, GIF, WebP"
	ErrInvalidDocumentType        = "invalid document type; allowed: PDF, DOC, DOCX"
	ErrUploadFailed               = "upload failed: "

	// Business logic errors
	ErrClinicCreatedButLinkFailed = "clinic was created but we could not link you as owner. Please try again."
	ErrInvalidJSONData            = "Invalid JSON data"
	ErrInvalidIDFormat            = "invalid id: "
)

// Success messages for HTTP responses
const (
	// Custom form messages
	MsgCustomFormCreated             = "custom form created"
	MsgCustomFormRetrieved           = "custom form retrieved"
	MsgCustomFormsRetrieved          = "custom forms retrieved"
	MsgPublishedCustomFormsRetrieved = "published custom forms retrieved"
	MsgCustomFormUpdated             = "custom form updated"
	MsgFormPublished                 = "form published"
	MsgFormUnpublished               = "form unpublished"
	MsgFormArchived                  = "form archived"
	MsgCustomFormDeleted             = "custom form deleted"
	MsgFormDuplicated                = "form duplicated"

	// Entry messages
	MsgEntryCreated          = "entry created"
	MsgFieldEntryRetrieved   = "field entry retrieved"
	MsgFieldEntriesRetrieved = "field entries retrieved"
	MsgNetDetailsRetrieved   = "net details retrieved"
	MsgGrossDetailsRetrieved = "gross details retrieved"
	MsgFieldEntryUpdated     = "field entry updated"
	MsgFieldEntryDeleted     = "field entry deleted"

	// Clinic messages
	MsgClinicCreatedSuccessfully    = "clinic created successfully"
	MsgClinicDeletedSuccessfully    = "clinic deleted successfully"
	MsgClinicsRetrievedSuccessfully = "clinics retrieved successfully"
	// User clinic messages
	MsgUserAssociatedWithClinic         = "user associated with clinic successfully"
	MsgUserClinicsRetrievedSuccessfully = "user clinics retrieved successfully"
	MsgClinicUsersRetrievedSuccessfully = "clinic users retrieved successfully"
	MsgUserRemovedFromClinic            = "user removed from clinic successfully"

	// Transaction messages
	MsgTransactionCreated    = "transaction created"
	MsgTransactionRetrieved  = "transaction retrieved"
	MsgTransactionsRetrieved = "transactions retrieved"
	MsgTransactionUpdated    = "transaction updated"
	MsgTransactionPosted     = "transaction posted"
	MsgTransactionVoided     = "transaction voided"
	MsgTransactionDeleted    = "transaction deleted"

	// P&L report messages
	MsgPnlReportGenerated   = "P&L report generated"
	MsgPnlReportRetrieved   = "P&L report retrieved"
	MsgPnlReportsRetrieved  = "P&L reports retrieved"
	MsgPnlReportFinalized   = "P&L report finalized"
	MsgPnlReportRegenerated = "P&L report regenerated"
	MsgPnlReportDeleted     = "P&L report deleted"

	// Auth messages
	MsgLoggedOutSuccessfully = "logged out successfully"
)
