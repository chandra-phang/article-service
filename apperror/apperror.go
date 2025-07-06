package apperror

import "errors"

var (
	// DB
	ErrObjectNotExists         = errors.New("object does not exist")
	ErrCreateRecordFailed      = errors.New("create record failed")
	ErrGetRecordFailed         = errors.New("get record failed")
	ErrScanRecordFailed        = errors.New("scan record failed")
	ErrStartTransactionFailed  = errors.New("start transaction failed")
	ErrCommitTransactionFailed = errors.New("commit transaction failed")
	ErrNoAffectedRows          = errors.New("no affected rows")

	// ElasticSearch
	ErrIndexElasticFailed  = errors.New("index to elastic failed")
	ErrSearchElasticFailed = errors.New("search on elastic failed")

	// Controller
	ErrUnmarshalRequestBodyFailed = errors.New("unmarshal request body failed")

	// Author
	ErrAuthorNotFound = errors.New("author not found")
)
