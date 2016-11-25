// +build !windows

package oleutil

import ole "github.com/nogio/noggo/depend/go-ole"

// ConnectObject creates a connection point between two services for communication.
func ConnectObject(disp *ole.IDispatch, iid *ole.GUID, idisp interface{}) (uint32, error) {
	return 0, ole.NewError(ole.E_NOTIMPL)
}
