package errorcode_test

import (
	"errors"
	"reflect"
	"testing"

	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

func TestFromCS3Status(t *testing.T) {
	var tests = []struct {
		status *cs3rpc.Status
		err    error
		ignore []cs3rpc.Code
		result *errorcode.Error
	}{
		{nil, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "unspecified error has occurred"))},
		{nil, errors.New("test error"), nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "test error"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_OK}, nil, nil, nil},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_NOT_FOUND}, nil, []cs3rpc.Code{cs3rpc.Code_CODE_NOT_FOUND}, nil},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_PERMISSION_DENIED}, nil, []cs3rpc.Code{cs3rpc.Code_CODE_NOT_FOUND, cs3rpc.Code_CODE_PERMISSION_DENIED}, nil},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_NOT_FOUND, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.ItemNotFound, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_PERMISSION_DENIED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.AccessDenied, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_UNAUTHENTICATED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.Unauthenticated, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_INVALID_ARGUMENT, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.InvalidRequest, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_ALREADY_EXISTS, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.NameAlreadyExists, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_FAILED_PRECONDITION, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.InvalidRequest, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_UNIMPLEMENTED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.NotSupported, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_INVALID, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_CANCELLED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_UNKNOWN, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_RESOURCE_EXHAUSTED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_ABORTED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_OUT_OF_RANGE, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.InvalidRange, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_INTERNAL, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_UNAVAILABLE, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.ServiceNotAvailable, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_REDIRECTION, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.GeneralException, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_INSUFFICIENT_STORAGE, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.QuotaLimitReached, "msg"))},
		{&cs3rpc.Status{Code: cs3rpc.Code_CODE_LOCKED, Message: "msg"}, nil, nil, conversions.ToPointer(errorcode.New(errorcode.ItemIsLocked, "msg"))},
	}

	for _, test := range tests {
		if output := errorcode.FromCS3Status(test.status, test.err, test.ignore...); !reflect.DeepEqual(output, test.result) {
			t.Error("Test Failed: {} expected, recieved: {}", test.result, output)
		}
	}
}

func TestFromStat(t *testing.T) {
	var tests = []struct {
		stat   *provider.StatResponse
		err    error
		result *errorcode.Error
	}{
		{nil, errors.New("some error"), conversions.ToPointer(errorcode.New(errorcode.GeneralException, "some error"))},
		{&provider.StatResponse{Status: &cs3rpc.Status{Code: cs3rpc.Code_CODE_OK}}, nil, nil},
	}

	for _, test := range tests {
		if output := errorcode.FromStat(test.stat, test.err); !reflect.DeepEqual(output, test.result) {
			t.Error("Test Failed: {} expected, recieved: {}", test.result, output)
		}
	}
}
