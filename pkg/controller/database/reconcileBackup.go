package database

import (
	"context"
	"fmt"
	kciv1alpha1 "github.com/kloeckner-i/db-operator/pkg/apis/kci/v1alpha1"
	backup "github.com/kloeckner-i/db-operator/pkg/controller/database/backup"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDatabase) createBackupJob(dbcr *kciv1alpha1.Database) error {
	if !dbcr.Spec.Backup.Enable {
		// if not enabled, skip
		return nil
	}

	cronjob, err := backup.GCSBackupCron(dbcr)
	if err != nil {
		return err
	}

	desiredCapacity, err := PersistentVolumeCapacity(dbcr)
	if err != nil {
		return fmt.Errorf("PersistentVolumeCapacity() failed - %s", err)
	}

	storageClassName, err := StorageClassName(dbcr)
	if err != nil {
		return fmt.Errorf("StorageClassName() failed - %s", err)
	}

	// (0) allocate a handle
	objectMeta := backup.ObjectMetaBuilder(dbcr, "volume")
	persistentVolumeFilesystem := v1.PersistentVolumeFilesystem

	// (1) specify a volume
	persistentVolume :=
		&v1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: objectMeta.Name, ResourceVersion: objectMeta.ResourceVersion}, Spec: v1.PersistentVolumeSpec{
			ClaimRef:               &v1.ObjectReference{Kind: "PersistentVolumeClaim", APIVersion: "v1", UID: objectMeta.UID, Namespace: objectMeta.Namespace, Name: objectMeta.Name},
			AccessModes:            []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Capacity:               v1.ResourceList{v1.ResourceName(v1.ResourceStorage): resource.MustParse(desiredCapacity)},
			PersistentVolumeSource: v1.PersistentVolumeSource{AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{VolumeID: objectMeta.Name}},

			StorageClassName: storageClassName, VolumeMode: &persistentVolumeFilesystem}}

	// (2) stake a claim
	persistentVolumeClaim := &v1.PersistentVolumeClaim{
		ObjectMeta: *objectMeta, Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: persistentVolume.Spec.AccessModes,
			Resources:   v1.ResourceRequirements{Requests: persistentVolume.Spec.Capacity},
			VolumeName:  objectMeta.Name,
			VolumeMode:  &persistentVolumeFilesystem,
		},
		Status: v1.PersistentVolumeClaimStatus{AccessModes: persistentVolume.Spec.AccessModes, Capacity: persistentVolume.Spec.Capacity}}

	// (3) create the claim
	pvc := &v1.PersistentVolumeClaim{}
	if err = controllerutil.SetControllerReference(dbcr, pvc, r.scheme); err != nil {
		return fmt.Errorf("DB: SetControllerReference failed - %s", err)
	}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: persistentVolumeClaim.Name, Namespace: persistentVolumeClaim.Namespace}, pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), persistentVolumeClaim)
			if err != nil {
				return fmt.Errorf("DB: Creating PVC failed, EINVAL - %s", err)
			}
		}
	} else {
		return fmt.Errorf("DB: Checking for existence of PVC failed, EINVAL - %s", err)
	}

	// (4) create the volume
	pv := &v1.PersistentVolume{}
	if err = controllerutil.SetControllerReference(dbcr, pv, r.scheme); err != nil {
		return fmt.Errorf("DB: SetControllerReference failed - %s", err)
	}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: persistentVolume.Name, Namespace: persistentVolume.Namespace}, pv)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), persistentVolume)
			if err != nil {
				return fmt.Errorf("DB: Creating PVC failed, EINVAL - %s", err)
			}
		}
	} else {
		return fmt.Errorf("DB: Checking for existence of PVC failed, EINVAL - %s", err)
	}

	// (5) create the job which uses the pv/pvc
	controllerutil.SetControllerReference(dbcr, cronjob, r.scheme)
	err = r.client.Create(context.TODO(), cronjob)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			// if resource already exists, update
			err = r.client.Update(context.TODO(), cronjob)
			if err != nil {
				logrus.Errorf("DB: namespace=%s, name=%s failed updating backup cronjob", dbcr.Namespace, dbcr.Name)
				return err
			}
		} else {
			// failed to create deployment
			logrus.Errorf("DB: namespace=%s, name=%s failed creating backup cronjob", dbcr.Namespace, dbcr.Name)
			return err
		}
	}

	return nil
}

func PersistentVolumeCapacity(dbcr *kciv1alpha1.Database) (string, error) {
	var capacity = ""

	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		return capacity, err
	}

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return capacity, fmt.Errorf("GetBackedType(() failed - %s", err)
		return capacity, err
	}

	switch backend {
	case "generic":
		return capacity, nil
	case "google":
		return capacity, nil
	case "amazon":

		if instance.Spec.Amazon.Capacity != "" {
			return instance.Spec.Amazon.Capacity, nil
		}

		if conf.Instances.Amazon.Capacity != "" {
			return conf.Instances.Amazon.Capacity, nil
		}
		return capacity, nil
	default:
		return capacity, fmt.Errorf("GetBackedType() returned unknown backup type of '%s'", backend)
	}
}

func StorageClassName(dbcr *kciv1alpha1.Database) (string, error) {
	var storageClassName = ""

	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		return storageClassName, err
	}

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return storageClassName, fmt.Errorf("GetBackedType(() failed - %s", err)
	}

	switch backend {
	case "generic":
		return storageClassName, nil
	case "google":
		return storageClassName, nil
	case "amazon":

		if instance.Spec.Amazon.StorageClassName != "" {
			return instance.Spec.Amazon.StorageClassName, nil
		}

		if conf.Instances.Amazon.StorageClassName != "" {
			return conf.Instances.Amazon.StorageClassName, nil
		}
		return storageClassName, nil
	default:
		return storageClassName, fmt.Errorf("GetBackedType() returned unknown backup type of '%s'", backend)
	}
}
