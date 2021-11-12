package go_request_id

import (
	"context"
	"regexp"
	"testing"
)

func TestExtractRequestID(t *testing.T) {
	txid := "abc"
	ctx := WithRequestID(context.Background(), txid)
	txid2 := ExtractRequestID(ctx)
	if txid2 != txid {
		t.Fatalf("expected txid %v, got %v", txid, txid2)
	}

	txid = ExtractRequestID(context.Background())
	if txid != "" {
		t.Fatalf("expected empty txid, got %v", txid)
	}

	// wrong type (not string)
	ctx = context.WithValue(context.Background(), keyRequestID, 42)
	txid = ExtractRequestID(ctx)
	if txid != "" {
		t.Fatalf("expected empty txid for wrong type, got %v (%[1]T)", txid)
	}
}

func TestWithNewRandomRequestID(t *testing.T) {
	ctx := WithNewRandomRequestID(context.Background())
	txid := ExtractRequestID(ctx)
	ok, err := regexp.MatchString(
		`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		txid,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("unexpected txid format: %v", txid)
	}
}
