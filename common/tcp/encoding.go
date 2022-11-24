package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	PkgBegin   = 0xEF
	PkgEnd     = 0xFF
	PkgVersion = 0x01
	PkgHeadLen = 11
)

type PkgHead struct {
	Version       byte   // 版本号
	MessageFlag   byte   // 区分消息类型
	Txid          uint32 // txid
	SizeOfMessage uint32 // 消息大小
}

// 封包
// bodylen=len(body)
// 0--------1--------2--------3--------7--------11--------X--------X+1
// | start  | version| flags  | txid   | bodylen | body   | end    |
// +--------+--------+--------+--------+---------+--------+--------+
func Pack(body []byte, messageFlag byte, txid uint32) []byte {
	bodyLen := uint32(len(body))

	buf := &bytes.Buffer{}
	buf.WriteByte(PkgBegin)
	buf.WriteByte(PkgVersion)
	buf.WriteByte(messageFlag)
	binary.Write(buf, binary.BigEndian, txid)
	binary.Write(buf, binary.BigEndian, bodyLen)
	buf.Write(body)
	buf.WriteByte(PkgEnd)

	return buf.Bytes()
}

// 解包 body messageflag txid err
func Unpack(reader io.Reader) ([]byte, byte, uint32, error) {
	headerData := make([]byte, PkgHeadLen)
	if _, err := io.ReadFull(reader, headerData); err != nil {
		return nil, 0, 0, err
	}

	headerBuf := bytes.NewBuffer(headerData)
	begin, _ := headerBuf.ReadByte()
	if begin != PkgBegin {
		return nil, 0, 0, fmt.Errorf("invalid pkghead:%v must be:%v", begin, PkgBegin)
	}

	var err error
	header := &PkgHead{}
	header.Version, err = headerBuf.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}

	header.MessageFlag, err = headerBuf.ReadByte()
	if err != nil {
		return nil, 0, 0, err
	}

	binary.Read(headerBuf, binary.BigEndian, &header.Txid)
	binary.Read(headerBuf, binary.BigEndian, &header.SizeOfMessage)

	data := make([]byte, header.SizeOfMessage+1)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, 0, 0, err
	}

	if data[len(data)-1] != PkgEnd {
		return nil, 0, 0, fmt.Errorf("invalid pkgend:%v must be:%v", data[len(data)-1], PkgEnd)
	}

	return data[:len(data)-1], header.MessageFlag, header.Txid, nil
}

func PackWrite(writer io.Writer, messageFlag byte, txid uint32, buff []byte) error {
	var size int

	data := Pack(buff, messageFlag, txid)
	for {
		body := data[size:]
		n, err := writer.Write(body)
		if err != nil {
			return err
		}

		size = size + n
		if size >= len(data) {
			break
		}
	}
	if size != len(data) {
		return fmt.Errorf("data len:%v actual send:%v", len(data), size)
	}

	return nil
}
