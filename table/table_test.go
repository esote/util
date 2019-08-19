package table

import (
	"encoding/hex"
	"testing"

	"github.com/esote/util/uuid"
)

const (
	key  = "9baacc8baed73d1f115d10d069a4ee63i"
	name = "test_table"
)

func TestTable(t *testing.T) {
	table, err := NewTable(name, 2, 2*uuid.LenUUID)

	if err != nil {
		t.Fatal(err)
	}

	const n = 10

	rows := make([]string, n)

	for i := 0; i < n; i++ {
		uuid, err := uuid.NewUUID()

		if err != nil {
			t.Fatal(err)
		}

		row := hex.EncodeToString(uuid)

		if err = table.Insert(key, row); err != nil {
			t.Fatal(err)
		}

		rows[i] = row
	}

	// Check IndexN matches expected list of messages.
	for take := 0; take <= n; take++ {
		index, err := table.IndexN(key, uint64(take))

		if err != nil {
			t.Fatal(err)
		}

		for i, m := range rows[len(rows)-take:] {
			if m != index[len(index)-i-1] {
				t.Fatalf("mismatch at take %d", take)
			}
		}
	}

	if err = table.Splay.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}

func TestTableUnique(t *testing.T) {
	table, err := NewTable(name, 2, 1)

	if err != nil {
		t.Fatal(err)
	}

	rows := []string{"a", "b", "c", "d"}

	insert := []string{"a", "b", "c", "b", "d", "a", "d"}

	for _, row := range insert {
		if err = table.InsertUnique(key, row); err != nil {
			t.Fatal(err)
		}
	}

	index, err := table.IndexN(key, 0)

	if err != nil {
		t.Fatal(err)
	}

	if len(index) != len(rows) {
		t.Fatal("incorrect length")
	}

	for i, row := range index {
		if rows[len(rows)-i-1] != row {
			t.Fatalf("row %d mismatch", i)
		}
	}

	if err = table.Splay.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}
