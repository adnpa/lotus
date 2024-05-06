package znet

import (
	"errors"
	"sync"

	"github.com/adnpa/lotus/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn
}

func (cm *ConnManager) RemoveConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
}

func (cm *ConnManager) GetConn(connId uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn not found")
	}
}

func (cm *ConnManager) GetConnNum() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearAllConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}
}
