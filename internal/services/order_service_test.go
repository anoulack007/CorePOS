package services

import (
	"errors"
	"math"
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
)

func TestNormalizeOrderLinesCombinesDuplicatesAndSorts(t *testing.T) {
	productA := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	productB := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	lines, err := normalizeOrderLines([]ports.OrderLineInput{
		{ProductID: productB, Quantity: 1},
		{ProductID: productA, Quantity: 2},
		{ProductID: productB, Quantity: 3},
	})
	if err != nil {
		t.Fatalf("normalizeOrderLines() error = %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 normalized lines, got %d", len(lines))
	}
	if lines[0].ProductID != productA || lines[0].Quantity != 2 {
		t.Fatalf("unexpected first line: %#v", lines[0])
	}
	if lines[1].ProductID != productB || lines[1].Quantity != 4 {
		t.Fatalf("unexpected second line: %#v", lines[1])
	}
}

func TestNormalizeOrderLinesRejectsInvalidInput(t *testing.T) {
	productID := uuid.New()
	tests := [][]ports.OrderLineInput{
		nil,
		{{ProductID: uuid.Nil, Quantity: 1}},
		{{ProductID: uuid.New(), Quantity: 0}},
		{{ProductID: uuid.New(), Quantity: -1}},
		{{ProductID: productID, Quantity: math.MaxInt}, {ProductID: productID, Quantity: 1}},
	}
	for _, lines := range tests {
		if _, err := normalizeOrderLines(lines); !errors.Is(err, domain.ErrInvalidOrder) {
			t.Fatalf("expected ErrInvalidOrder for %#v, got %v", lines, err)
		}
	}
}

func TestMoneyUsesCentRounding(t *testing.T) {
	if got := moneyToCents(19.99); got != 1999 {
		t.Fatalf("moneyToCents(19.99) = %d", got)
	}
	if got := moneyToCents(0.1 + 0.2); got != 30 {
		t.Fatalf("moneyToCents(0.1 + 0.2) = %d", got)
	}
	if got := centsToMoney(1999); got != 19.99 {
		t.Fatalf("centsToMoney(1999) = %v", got)
	}
}
