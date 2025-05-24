// server/server.go
package models

type ServerMetaData interface {
	GetID() string
	GetHashVal() int
	GetContainerID() string
	GetContainerHostName() string
	GetServerID() string

	SetContainerID(containerID string)
	SetHashVal(hashVal int)
	SetID(id string)
}

type ServerInstance struct {
	ID          string
	HashVal     int
	containerID string
}

// Getter Methods
func (s *ServerInstance) GetID() string {
	return s.ID
}
func (s *ServerInstance) GetHashVal() int {
	return s.HashVal
}
func (s *ServerInstance) GetContainerID() string {
	return s.containerID
}
func (s *ServerInstance) GetContainerHostName() string {
	return "serverinst_" + s.ID
}
func (s *ServerInstance) GetServerID() string {
	return "Server_" + s.ID
}

// Setter Methods
func (s *ServerInstance) SetContainerID(containerID string) {
	s.containerID = containerID
}
func (s *ServerInstance) SetHashVal(hashVal int) {
	s.HashVal = hashVal
}
func (s *ServerInstance) SetID(id string) {
	s.ID = id
}
