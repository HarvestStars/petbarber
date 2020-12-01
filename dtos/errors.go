package dtos

const (
	OK             = 0
	UNKNOW_REQUEST = 1

	IMAGE_FETCH_ERROR  = 100
	IMAGE_UPLOAD_ERROR = 101

	LOGIN_SMS_CODE_MISSMATCH    = 10001
	LOGIN_SMS_CODE_INVALID      = 10002
	LOGIN_SMS_CODE_EXPIRED      = 10003
	LOGIN_UPGRADE_INVALID_TOKEN = 10004

	JWT_VERIFY_RESULT_EXPIRED     = 20001
	JWT_VERIFY_RESULT_BAD_TOKEN   = 20002
	JWT_EXPECTED_PETHOUSE_TOKEN   = 20003
	JWT_EXPECTED_PETGROOMER_TOKEN = 20004
	JWT_MISSING_TOKEN             = 20005
	JWT_CREATE_WRONG              = 20006
	JWT_TYPE_WRONG                = 20007

	ORDER_UNKNOWN_ORDER_TYPE   = 30001
	ORDER_TYPE_BODY_MISMATCH   = 30002
	ORDER_PAYMENT_DATA_MISSION = 30003
	ORDER_NOT_EXISTS           = 30004
	ORDER_OCCUPIED             = 30005
	ORDER_DENIED_USER          = 30006
	ORDER_INTERNAL_ERROR       = 30007
	ORDER_NOT_ACTIVE           = 30008
	ORDER_NOT_ACTIVE2          = 30009
	ORDER_CANCEL_NOT_ALLOWED   = 30010
	ORDER_NOT_FINISHED         = 30011

	COMMENT_ERROR_TYPE          = 31001
	COMMENT_CANT_CREATE_COMMENT = 31002
)