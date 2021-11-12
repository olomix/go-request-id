package go_request_id

import (
	"context"
	"regexp"
	"testing"
)

func TestExtractTxid(t *testing.T) {
	txid := "abc"
	ctx := WithTxid(context.Background(), txid)
	txid2 := ExtractTxid(ctx)
	if txid2 != txid {
		t.Fatalf("expected txid %v, got %v", txid, txid2)
	}

	txid = ExtractTxid(context.Background())
	if txid != "" {
		t.Fatalf("expected empty txid, got %v", txid)
	}

	// wrong type (not string)
	ctx = context.WithValue(context.Background(), keyTxid, 42)
	txid = ExtractTxid(ctx)
	if txid != "" {
		t.Fatalf("expected empty txid for wrong type, got %v (%[1]T)", txid)
	}
}

func TestWithNewRandomTxid(t *testing.T) {
	ctx := WithNewRandomTxid(context.Background())
	txid := ExtractTxid(ctx)
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
