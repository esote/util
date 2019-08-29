package table

import (
	"encoding/hex"
	"math/rand"
	"os"
	"testing"

	"github.com/esote/util/uuid"
)

const name = "test_table"

var key string

func TestMain(m *testing.M) {
	if table, err := NewTable(name, 1, 1); err == nil {
		_ = table.Splay.RemoveAll()
	}

	u, _ := uuid.NewUUID()
	key = hex.EncodeToString(u)

	os.Exit(m.Run())
}

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

func TestUnique(t *testing.T) {
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

func TestDelete(t *testing.T) {
	iterations := []struct {
		before []string
		expect []string
		remove string
	}{
		{
			// Left edge.
			before: []string{"a", "b"},
			expect: []string{"b"},
			remove: "a",
		},
		{
			// Right edge.
			before: []string{"a", "b"},
			expect: []string{"a"},
			remove: "b",
		},
		{
			// Middle.
			before: []string{"a", "b", "c", "d"},
			expect: []string{"d", "c", "a"},
			remove: "b",
		},
		{
			// Only.
			before: []string{"abc"},
			expect: []string{},
			remove: "abc",
		},
		{
			// Nonexistent.
			before: []string{"a", "b", "c"},
			expect: []string{"c", "b", "a"},
			remove: "d",
		},
	}

	for _, it := range iterations {
		table, err := NewTable(name, 2, len(it.remove))

		if err != nil {
			t.Fatal(err)
		}

		for _, row := range it.before {
			if err = table.Insert(key, row); err != nil {
				t.Fatal(err)
			}
		}

		if err = table.Delete(key, it.remove); err != nil {
			t.Fatal(err)
		}

		index, err := table.IndexN(key, 0)

		if err != nil {
			t.Fatal(err)
		}

		if len(index) != len(it.expect) {
			t.Fatal("invalid length")
		}

		for i, row := range index {
			if row != it.expect[i] {
				t.Fatalf("row %d mismatch", i)
			}
		}

		if err = table.Splay.RemoveAll(); err != nil {
			t.Fatal(err)
		}
	}

}

const (
	bufsize = 100
	rowsize = 10
)

// BenchmarkDelete100 benchmarks deleting all rows from a table of 100 rows.
func BenchmarkDelete100(b *testing.B) {
	b.StopTimer()

	var buf [bufsize]string

	for i := range buf {
		b := make([]byte, rowsize)
		rand.Read(b)
		buf[i] = string(b)
	}

	table, _ := NewTable(name, 2, rowsize)

	for i := 0; i < b.N; i++ {
		for i := range buf {
			_ = table.Insert(key, buf[i])
		}

		b.StartTimer()
		for row := range buf {
			_ = table.Delete(key, buf[row])
		}
		b.StopTimer()

		_ = table.Splay.Remove(key)

	}

	if err := table.Splay.RemoveAll(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkIndexN100 benchmarks indexing all rows from a table of 100 rows.
func BenchmarkIndexN100(b *testing.B) {
	b.StopTimer()

	var buf [bufsize]string

	for i := range buf {
		b := make([]byte, rowsize)
		rand.Read(b)
		buf[i] = string(b)
	}

	table, _ := NewTable(name, 2, rowsize)

	for i := 0; i < b.N; i++ {
		for i := range buf {
			_ = table.Insert(key, buf[i])
		}

		b.StartTimer()
		_, _ = table.IndexN(key, 0)
		b.StopTimer()

		_ = table.Splay.Remove(key)

	}

	if err := table.Splay.RemoveAll(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkInsert100 benchmarks inserting 100 rows into an empty table.
func BenchmarkInsert100(b *testing.B) {
	b.StopTimer()

	var buf [bufsize]string

	for i := range buf {
		b := make([]byte, rowsize)
		rand.Read(b)
		buf[i] = string(b)
	}

	table, _ := NewTable(name, 2, rowsize)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for i := range buf {
			_ = table.Insert(key, buf[i])
		}
		b.StopTimer()

		_ = table.Splay.Remove(key)
	}

	if err := table.Splay.RemoveAll(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkInsertUnique100 benchmarks inserting 100 rows into an empty table where
// 10 percent of the inserts are duplicates.
func BenchmarkInsertUnique100(b *testing.B) {
	b.StopTimer()

	var buf [bufsize]string

	for i := range buf {
		b := make([]byte, rowsize)
		rand.Read(b)
		buf[i] = string(b)
	}

	// Introduce duplicates.
	for i := 0; i < len(buf); i += len(buf) / 10 {
		buf[i] = buf[rand.Intn(len(buf))]
	}

	table, _ := NewTable(name, 2, rowsize)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		for row := range buf {
			_ = table.InsertUnique(key, buf[row])
		}
		b.StopTimer()

		_ = table.Splay.Remove(key)
	}

	if err := table.Splay.RemoveAll(); err != nil {
		b.Fatal(err)
	}
}
