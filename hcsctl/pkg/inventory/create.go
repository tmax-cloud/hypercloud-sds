package inventory

import (
	"context"

	"github.com/codeclysm/extract"
	"github.com/golang/glog"
	"github.com/markbates/pkger"

	"os"
	"path"
	"path/filepath"
)

const (
	sampleInventoryPrefix = "/sample-inventory"
	tarName               = "hcs.tar.gz"
)

// Create creates inventory with fixed version
// Usage: Create(inventoryName)
func Create(inventoryName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	inventoryDir := path.Join(wd, inventoryName)
	glog.Infof("Start Creating Sample Inventory on \"%s\"", inventoryDir)
	glog.Infof("Sample Inventory contains rook version : \"v1.3.6\"") // TODO yaml 파일 읽어서 값 가져오는 것으로 변경
	glog.Infof("Sample Inventory contains cdi version : \"v1.18.0\"")

	err = createInventory(inventoryDir)

	return err
}

func createInventory(inventory string) error {
	err := os.MkdirAll(inventory, 0755)
	if err != nil {
		return err
	}

	includedTarPath := filepath.Join(sampleInventoryPrefix, tarName)
	sourceTarFile, err := pkger.Open(includedTarPath)

	if err != nil {
		return err
	}

	defer func() {
		err = sourceTarFile.Close()
		if err != nil {
			glog.Fatalln(err)
			os.Exit(1)
		}
	}()

	err = extract.Gz(context.TODO(), sourceTarFile, inventory, nil)
	if err != nil {
		return err
	}

	return nil
}
