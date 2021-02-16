package conv

import (
	"testing"
)

func TestInt2Letter(t *testing.T) {
	if Int2Letter(1) != "A" {
		t.Errorf("except A, but %s", Int2Letter(1))
	}
	if Int2Letter(25) != "Y" {
		t.Errorf("except Y, but %s", Int2Letter(25))
	}
	if Int2Letter(26) != "Z" {
		t.Errorf("except Z, but %s", Int2Letter(26))
	}
	if Int2Letter(27) != "AA" {
		t.Errorf("except AA, but %s", Int2Letter(27))
	}
	if Int2Letter(100) != "CV" {
		t.Errorf("except CV, but %s", Int2Letter(100))
	}
}