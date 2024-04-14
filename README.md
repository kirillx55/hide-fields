# Hide fields
Hide fields in go stucture by tag 'hide'

## Supported types
int, int8, int16, int32, int64,
uint, uint8, uint16, uint32, uint64,
float32, float64, complex64, complex128,
bool, string

array, slice, map

#### *Please note:*
* *that complex64 will hide by default type value only*
* *not hide if field value equal default type value*
* *be careful with pointers cause this method mutate structure values*
* *Doesn't support interfaces*

### how to
```go
type t struct {
	Int         int               `hide:"-1"`
	Uint        uint              `hide:"-1"`
	str         string            `hide:"**"`
	Slice       []string          `hide:"**"`
	Map         map[string]string `hide:"**"`
}

v := t{
	 Int:   1,
	 Uint:  2,
	 str:   "string",
	 Slice: []string{"one", "", "three"},
	 Map:   map[string]string{"key": "value"},
}

HideFields(&v)
```

```go
// the resust is v with hid values
{
Int:   -1,
Uint:  0, // cause "-1" is not valid value for this type 
str:   "**",
Slice: []string{"**", "", "**"},
Map:   map[string]string{"key": "**"},
}
```