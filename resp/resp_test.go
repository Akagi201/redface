package resp_test

import (
	"bytes"
	"errors"
	"fmt"
	. "testing"

	"github.com/Akagi201/redface/resp"
	"github.com/stretchr/testify/assert"
)

func pretendRead(s string) *resp.Resp {
	buf := bytes.NewBufferString(s)
	return resp.NewRespReader(buf).Read()
}

func TestRead(t *T) {

	r := pretendRead("")
	assert.NotNil(t, r.Err)

	// Simple string
	r = pretendRead("+ohey\r\n")
	assert.Equal(t, resp.SimpleStr, r.Typ)
	assert.Exactly(t, []byte("ohey"), r.Val)
	s, err := r.Str()
	assert.Nil(t, err)
	assert.Equal(t, "ohey", s)

	// Empty simple string
	r = pretendRead("+\r\n")
	assert.Equal(t, resp.SimpleStr, r.Typ)
	assert.Exactly(t, []byte(""), r.Val)
	s, err = r.Str()
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	// Error
	r = pretendRead("-ohey\r\n")
	assert.Equal(t, resp.AppErr, r.Typ)
	assert.Exactly(t, errors.New("ohey"), r.Val)
	assert.Equal(t, "ohey", r.Err.Error())

	// Empty error
	r = pretendRead("-\r\n")
	assert.Equal(t, resp.AppErr, r.Typ)
	assert.Exactly(t, errors.New(""), r.Val)
	assert.Equal(t, "", r.Err.Error())

	// Int
	r = pretendRead(":1024\r\n")
	assert.Equal(t, resp.Int, r.Typ)
	assert.Exactly(t, int64(1024), r.Val)
	i, err := r.Int()
	assert.Nil(t, err)
	assert.Equal(t, 1024, i)

	// Int (from string)
	r = pretendRead("+50\r\n")
	assert.Equal(t, resp.SimpleStr, r.Typ)
	assert.Exactly(t, []byte("50"), r.Val)
	i, err = r.Int()
	assert.Nil(t, err)
	assert.Equal(t, 50, i)

	// Int (from string, can't parse)
	r = pretendRead("+ImADuck\r\n")
	assert.Equal(t, resp.SimpleStr, r.Typ)
	assert.Exactly(t, []byte("ImADuck"), r.Val)
	_, err = r.Int()
	assert.NotNil(t, err)

	// Bulk string
	r = pretendRead("$3\r\nfoo\r\n")
	assert.Equal(t, resp.BulkStr, r.Typ)
	assert.Exactly(t, []byte("foo"), r.Val)
	s, err = r.Str()
	assert.Nil(t, err)
	assert.Equal(t, "foo", s)

	// Empty bulk string
	r = pretendRead("$0\r\n\r\n")
	assert.Equal(t, resp.BulkStr, r.Typ)
	assert.Exactly(t, []byte(""), r.Val)
	s, err = r.Str()
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	// Nil bulk string
	r = pretendRead("$-1\r\n")
	assert.Equal(t, resp.Nil, r.Typ)

	// Array
	r = pretendRead("*2\r\n+foo\r\n+bar\r\n")
	assert.Equal(t, resp.Array, r.Typ)
	assert.Equal(t, 2, len(r.Val.([]resp.Resp)))
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[0].Typ)
	assert.Exactly(t, []byte("foo"), r.Val.([]resp.Resp)[0].Val)
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[1].Typ)
	assert.Exactly(t, []byte("bar"), r.Val.([]resp.Resp)[1].Val)
	l, err := r.List()
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "bar"}, l)
	b, err := r.ListBytes()
	assert.Nil(t, err)
	assert.Equal(t, [][]byte{[]byte("foo"), []byte("bar")}, b)
	m, err := r.Map()
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"foo": "bar"}, m)

	// Empty Array
	r = pretendRead("*0\r\n")
	assert.Equal(t, resp.Array, r.Typ)
	assert.Equal(t, 0, len(r.Val.([]resp.Resp)))

	// Nil Array
	r = pretendRead("*-1\r\n")
	assert.Equal(t, resp.Nil, r.Typ)

	// Embedded Array
	r = pretendRead("*3\r\n+foo\r\n+bar\r\n*2\r\n+foo\r\n+bar\r\n")
	assert.Equal(t, resp.Array, r.Typ)
	assert.Equal(t, 3, len(r.Val.([]resp.Resp)))
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[0].Typ)
	assert.Exactly(t, []byte("foo"), r.Val.([]resp.Resp)[0].Val)
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[1].Typ)
	assert.Exactly(t, []byte("bar"), r.Val.([]resp.Resp)[1].Val)
	r = &r.Val.([]resp.Resp)[2]
	assert.Equal(t, 2, len(r.Val.([]resp.Resp)))
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[0].Typ)
	assert.Exactly(t, []byte("foo"), r.Val.([]resp.Resp)[0].Val)
	assert.Equal(t, resp.SimpleStr, r.Val.([]resp.Resp)[1].Typ)
	assert.Exactly(t, []byte("bar"), r.Val.([]resp.Resp)[1].Val)

	// Test that two bulks in a row read correctly
	r = pretendRead("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	assert.Equal(t, resp.Array, r.Typ)
	assert.Equal(t, 2, len(r.Val.([]resp.Resp)))
	assert.Equal(t, resp.BulkStr, r.Val.([]resp.Resp)[0].Typ)
	assert.Exactly(t, []byte("foo"), r.Val.([]resp.Resp)[0].Val)
	assert.Equal(t, resp.BulkStr, r.Val.([]resp.Resp)[1].Typ)
	assert.Exactly(t, []byte("bar"), r.Val.([]resp.Resp)[1].Val)
}

type arbitraryTest struct {
	val    interface{}
	expect []byte
}

var nilResp = pretendRead("$-1\r\n")

var arbitraryTests = []arbitraryTest{
	{[]byte("OHAI"), []byte("$4\r\nOHAI\r\n")},
	{"OHAI", []byte("$4\r\nOHAI\r\n")},
	{true, []byte("$1\r\n1\r\n")},
	{false, []byte("$1\r\n0\r\n")},
	{nil, []byte("$-1\r\n")},
	{80, []byte(":80\r\n")},
	{int64(-80), []byte(":-80\r\n")},
	{uint64(80), []byte(":80\r\n")},
	{float32(0.1234), []byte("$6\r\n0.1234\r\n")},
	{float64(0.1234), []byte("$6\r\n0.1234\r\n")},
	{errors.New("hi"), []byte("-hi\r\n")},

	{nilResp, []byte("$-1\r\n")},

	{[]int{1, 2, 3}, []byte("*3\r\n:1\r\n:2\r\n:3\r\n")},
	{map[int]int{1: 2}, []byte("*2\r\n:1\r\n:2\r\n")},

	{resp.NewRespSimple("OK"), []byte("+OK\r\n")},
}

var arbitraryAsFlattenedStringsTests = []arbitraryTest{
	{
		[]interface{}{"wat", map[string]interface{}{
			"foo": 1,
		}},
		[]byte("*3\r\n$3\r\nwat\r\n$3\r\nfoo\r\n$1\r\n1\r\n"),
	},
	{map[string]interface{}{"foo": true}, []byte("*2\r\n$3\r\nfoo\r\n$1\r\n1\r\n")},
}

func TestWriteArbitrary(t *T) {
	var err error
	buf := bytes.NewBuffer([]byte{})
	for _, test := range arbitraryTests {
		buf.Reset()
		_, err = resp.NewResp(test.val).WriteTo(buf)
		assert.Nil(t, err)
		assert.Equal(t, test.expect, buf.Bytes())
	}
}

func TestWriteArbitraryAsFlattenedStrings(t *T) {
	var err error
	buf := bytes.NewBuffer([]byte{})
	for _, test := range arbitraryAsFlattenedStringsTests {
		buf.Reset()
		_, err = resp.NewRespFlattenedStrings(test.val).WriteTo(buf)
		assert.Nil(t, err)
		assert.Equal(t, test.expect, buf.Bytes())
	}
}

func TestFloat64(t *T) {
	r := resp.NewResp(4)
	_, err := r.Float64()
	assert.NotNil(t, err)

	testErr := fmt.Errorf("test")
	r = resp.NewResp(testErr)
	_, err = r.Float64()
	assert.Equal(t, testErr, err)

	r = resp.NewResp("test")
	_, err = r.Float64()
	assert.NotNil(t, err)

	r = resp.NewResp("5.0")
	f, err := r.Float64()
	assert.Nil(t, err)
	assert.Equal(t, float64(5.0), f)
}
