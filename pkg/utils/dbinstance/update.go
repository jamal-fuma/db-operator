package dbinstance

import (
	"github.com/sirupsen/logrus"
)

// Update instance if instance exists
func Update(ins DbInstance) (map[string]string, error) {
	err := ins.exist()
	if err != nil {
		return nil, ErrNotExists
	}

	logrus.Debug("instance exists, update instance")
	err = ins.update()
	if err != nil {
		logrus.Debug("update failed")
		return nil, err
	}

	data, err := ins.getInfoMap()
	if err != nil {
		return nil, err
	}

	return data, nil
}
