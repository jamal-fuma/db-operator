package dbinstance

// Amazon represents database instance which can be connected by address and port
type Amazon struct {
	Generic            Generic
	ServiceAccountName string
	Capacity           string
	StorageClassName   string
	FSGroup            int64
}

func (ins *Amazon) exist() error {
	return ins.Generic.exist()
}

func (ins *Amazon) create() error {
	return ins.Generic.create()
}

func (ins *Amazon) update() error {
	return ins.Generic.update()
}

func (ins *Amazon) getInfoMap() (map[string]string, error) {
	return ins.Generic.getInfoMap()
}
