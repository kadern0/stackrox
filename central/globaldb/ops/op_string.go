// Code generated by "stringer -type=Op"; DO NOT EDIT.

package ops

import "strconv"

const _Op_name = "AddCountGetGetAllGetManyListRenameRemoveUpdateUpsert"

var _Op_index = [...]uint8{0, 3, 8, 11, 17, 24, 28, 34, 40, 46, 52}

func (i Op) String() string {
	if i < 0 || i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
