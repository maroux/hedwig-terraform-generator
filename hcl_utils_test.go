package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHclFmt(t *testing.T) {
	contents := []byte(`
variable "aws_account_id" {
}
`)
	goodContents := []byte(`variable "aws_account_id" {}
`)
	filename := "bad_hcl.tf"
	require.NoError(t, ioutil.WriteFile(filename, contents, os.ModePerm))
	defer os.Remove(filename)

	hclFmt(filename)

	contents, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	assert.Equal(t, string(goodContents), string(contents))
}

func TestHclFmtDir(t *testing.T) {
	contents := []byte(`
variable "aws_account_id" {
}
`)
	goodContents := []byte(`variable "aws_account_id" {}
`)
	filename := "bad_hcl.tf"
	require.NoError(t, ioutil.WriteFile(filename, contents, os.ModePerm))
	defer os.Remove(filename)

	hclFmtDir(".")

	contents, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	assert.Equal(t, string(goodContents), string(contents))
}

func TestHclValue_Empty(t *testing.T) {
	var emptySlice []string
	assert.Equal(t, "[]", hclvalue(emptySlice))
}

func TestHclValue_Nested(t *testing.T) {
	nestedList := []string{"sub", "list"}
	nestedMap := map[string]string{
		`mic"key`: "mouse",
		"foo":     "bar",
	}
	slice := []interface{}{1, nestedList, nestedMap}
	assert.Equal(t, `[
1,
[
"sub",
"list"
],
{
foo = "bar"
"mic\"key" = "mouse"
}
]`, hclvalue(slice))
}

func TestHclValue_Quote(t *testing.T) {
	value := `foo"bar`
	assert.Equal(t, `"foo\"bar"`, hclvalue(value))
}

func TestHclIdent(t *testing.T) {
	assert.Equal(t, `my-function`, hclident("myFunction"))
	assert.Equal(t, `my-function`, hclident("my_function"))
	// known limitation:
	assert.Equal(t, `${var.my-function}`, hclident("${var.myFunction}"))
}
