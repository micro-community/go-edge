package iobuffer

import "bytes"

//ReadWriteCloser ...
type ReadWriteCloser struct {
	wbuf *bytes.Buffer
	rbuf *bytes.Buffer
}

//NewBuffer return readWriteCloser
func NewBuffer() *ReadWriteCloser {

	return &ReadWriteCloser{
		rbuf: bytes.NewBuffer(nil),
		wbuf: bytes.NewBuffer(nil),
	}

}

func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.rbuf.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.wbuf.Write(p)
}

//WriteRbuf write ReaderBuffer
func (rwc *ReadWriteCloser) WriteRbuf(p []byte) (n int, err error) {
	return rwc.rbuf.Write(p)
}

//Close ...
func (rwc *ReadWriteCloser) Close() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()
	return nil
}

//Reset ...
func (rwc *ReadWriteCloser) Reset() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()
	return nil
}

//String ...
func (rwc *ReadWriteCloser) String() string {
	return "ReadWriteCloser"
}

//WBytes ...
func (rwc *ReadWriteCloser) WBytes() []byte {
	return rwc.wbuf.Bytes()
}

//RBytes ...
func (rwc *ReadWriteCloser) RBytes() []byte {
	return rwc.rbuf.Bytes()
}
