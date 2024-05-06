package zpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/zconf"
	"zinx/ziface"
)

type Datapack struct {
	HeadLen uint32
}

func NewDatapack() *Datapack {
	return &Datapack{}
}

// 包头8字节
func (dp *Datapack) GetHeadLen() uint32 {
	// DataLen uint32 | Id uint32
	return 8
}

func (dp *Datapack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})
	// 写入长度，id，数据，需要保证顺序
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

func (dp *Datapack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断datalen是否超出最大允许长度
	if zconf.GGlobalObj.MaxPackageSize > 0 && msg.DataLen > zconf.GGlobalObj.MaxPackageSize {
		return nil, errors.New("too Large msg data recv")
	}

	return msg, nil

}
