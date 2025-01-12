package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerSkipsIfNoRecordMatches(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(9999, 12, 31))
	require.Nil(t, reconciler)
}

func TestReconcilerRespectsIndentationStyle(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"1444-10-09\n\t1h", "1444-10-09\n\t1h\n\t30m\n"},
		{"1444-10-09\n  1h", "1444-10-09\n  1h\n  30m\n"},
		{"1444-10-09\n   1h", "1444-10-09\n   1h\n   30m\n"},
		{"1444-10-09\n    1h", "1444-10-09\n    1h\n    30m\n"},
		{"1444-10-08\n  3h\n\n1444-10-09\n\t1h", "1444-10-08\n  3h\n\n1444-10-09\n\t1h\n\t30m\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(1444, 10, 9))
		result, err := reconciler.AppendEntry("30m")
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcilerRespectsLineEndingStyle(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"1444-10-09\r\n\t1h", "1444-10-09\r\n\t1h\r\n\t30m\r\n"},
		{"1444-10-09\r\n\t1h\n\t2h", "1444-10-09\r\n\t1h\n\t2h\r\n\t30m\r\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(1444, 10, 9))
		result, err := reconciler.AppendEntry("30m")
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAddRecordIfOriginalIsEmpty(t *testing.T) {
	rs, _ := parser.Parse("")
	reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2000, 5, 5), nil)
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, "2000-05-05\n", result.AllSerialised)
	assert.Equal(t, "2000-05-05", result.Record.Date().ToString())
}

func TestReconcileAddRecordIfOriginalContainsOneRecord(t *testing.T) {
	rs, _ := parser.Parse("1999-12-31")
	reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2000, 2, 1), nil)
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, "1999-12-31\n\n2000-02-01\n", result.AllSerialised)
	assert.Equal(t, "2000-02-01", result.Record.Date().ToString())
}

func TestReconcileNewRecordFromEmptyFile(t *testing.T) {
	for _, x := range []struct {
		original string
	}{
		{""},
		{"\n"},
		{"\n\n"},
		{"\n\n\t\n"},
		{"\n\n     \t\n \t     \n  "},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(1995, 3, 17), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, "1995-03-17\n", result.AllSerialised)
	}
}

func TestReconcilePrependNewRecord(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-02", "2018-01-01\n\n2018-01-02"},
		{"2018-01-02", "2018-01-01\n\n2018-01-02"},
		{"2018-01-02\n\t1h", "2018-01-01\n\n2018-01-02\n\t1h"},
		{"2018-01-02\n\n", "2018-01-01\n\n2018-01-02\n\n"},
		{"\n\n2018-01-02\n", "2018-01-01\n\n\n\n2018-01-02\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2018, 1, 1), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAppendNewRecord(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-01", "2018-01-01\n\n2019-01-01\n"},
		{"2018-01-01\n\n", "2018-01-01\n\n2019-01-01\n\n"},
		{"\n\n2018-01-01\n", "\n\n2018-01-01\n\n2019-01-01\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2019, 1, 1), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAddBlockInBetween(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-01\n\n2018-01-03", "2018-01-01\n\n2018-01-02\n\n2018-01-03"},
		{"2018-01-01\n\n\n2018-01-03", "2018-01-01\n\n2018-01-02\n\n\n2018-01-03"},
		{"2018-01-02\n\t1h\n\n2018-01-03", "2018-01-02\n\t1h\n\n2018-01-02\n\n2018-01-03"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2018, 1, 2), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAddRecordWithShouldTotal(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(2018, 1, 2), NewShouldTotal(5, 31))
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02 (5h31m!)
`, result.AllSerialised)
	assert.Equal(t, NewShouldTotal(5, 31), result.Record.ShouldTotal())
}

func TestReconcileRespectsExistingStylePref(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"3145/06/15\n", "3145/06/15\n\n3145/06/16\n"},
		{"3145/06/14\n\n3145/06/15\n\n3145-06-15\n", "3145/06/14\n\n3145/06/15\n\n3145-06-15\n\n3145/06/16\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerForNewRecord(rs, Ɀ_Date_(3145, 6, 16), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}
