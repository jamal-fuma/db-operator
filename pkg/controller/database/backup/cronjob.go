package backup

import (
	"errors"
	"fmt"

	kciv1alpha1 "github.com/kloeckner-i/db-operator/pkg/apis/kci/v1alpha1"
	"github.com/kloeckner-i/db-operator/pkg/config"
	"github.com/kloeckner-i/db-operator/pkg/utils/kci"

	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

var conf = config.Config{}

// GCSBackupCron builds kubernetes cronjob object
// to create database backup regularly with defined schedule from dbcr
// this job will database dump and upload to google bucket storage for backup
func GCSBackupCron(dbcr *kciv1alpha1.Database) (*batchv1beta1.CronJob, error) {
	cronJobSpec, err := buildCronJobSpec(dbcr)
	if err != nil {
		return nil, err
	}

	return &batchv1beta1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch",
		},
		ObjectMeta: *ObjectMetaBuilder(dbcr, "backup"),
		Spec:       cronJobSpec,
	}, nil
}

func buildCronJobSpec(dbcr *kciv1alpha1.Database) (batchv1beta1.CronJobSpec, error) {
	jobTemplate, err := buildJobTemplate(dbcr)
	if err != nil {
		return batchv1beta1.CronJobSpec{}, err
	}

	return batchv1beta1.CronJobSpec{
		JobTemplate: jobTemplate,
		Schedule:    dbcr.Spec.Backup.Cron,
	}, nil
}

func buildJobTemplate(dbcr *kciv1alpha1.Database) (batchv1beta1.JobTemplateSpec, error) {
	ActiveDeadlineSeconds := int64(60 * 10) // 10m
	BackoffLimit := int32(3)
	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		logrus.Errorf("can not build job template - %s", err)
		return batchv1beta1.JobTemplateSpec{}, err
	}

	account, err := getServiceAccountName(dbcr)
	if err != nil {
		logrus.Errorf("can not build job template - %s", err)
		return batchv1beta1.JobTemplateSpec{}, err
	}

	securityContext, err := getSecurityContext(dbcr)
	if err != nil {
		logrus.Errorf("can not build job template - %s", err)
		return batchv1beta1.JobTemplateSpec{}, err
	}

	var backupContainer v1.Container

	engine := instance.Spec.Engine
	switch engine {
	case "postgres":
		backupContainer, err = postgresBackupContainer(dbcr)
		if err != nil {
			return batchv1beta1.JobTemplateSpec{}, err
		}
	case "mysql":
		backupContainer, err = mysqlBackupContainer(dbcr)
		if err != nil {
			return batchv1beta1.JobTemplateSpec{}, err
		}
	default:
		return batchv1beta1.JobTemplateSpec{}, errors.New("unknown engine type")
	}

	vols, err := volumes(dbcr)
	if err != nil {
		return batchv1beta1.JobTemplateSpec{}, err
	}

	return batchv1beta1.JobTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: kci.BaseLabelBuilder(),
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: &ActiveDeadlineSeconds,
			BackoffLimit:          &BackoffLimit,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: kci.BaseLabelBuilder(),
				},
				Spec: v1.PodSpec{
					Containers:         []v1.Container{backupContainer},
					NodeSelector:       conf.Backup.NodeSelector,
					ServiceAccountName: account,
					RestartPolicy:      v1.RestartPolicyNever,
					Volumes:            vols,
					SecurityContext:    securityContext,
				},
			},
		},
	}, nil
}

func postgresBackupContainer(dbcr *kciv1alpha1.Database) (v1.Container, error) {
	env, err := postgresEnvVars(dbcr)
	if err != nil {
		return v1.Container{}, err
	}
	mounts, err := volumeMounts(dbcr)
	if err != nil {
		return v1.Container{}, err
	}

	return v1.Container{
		Name:            "postgres-dump",
		Image:           conf.Backup.Postgres.Image,
		ImagePullPolicy: v1.PullAlways,
		VolumeMounts:    mounts,
		Env:             env,
	}, nil
}

func mysqlBackupContainer(dbcr *kciv1alpha1.Database) (v1.Container, error) {
	env, err := mysqlEnvVars(dbcr)
	if err != nil {
		return v1.Container{}, err
	}
	mounts, err := volumeMounts(dbcr)
	if err != nil {
		return v1.Container{}, err
	}

	return v1.Container{
		Name:            "mysql-dump",
		Image:           conf.Backup.Mysql.Image,
		ImagePullPolicy: v1.PullAlways,
		VolumeMounts:    mounts,
		Env:             env,
	}, nil
}

func volumeMounts(dbcr *kciv1alpha1.Database) ([]v1.VolumeMount, error) {
	mounts := make([]v1.VolumeMount, 0)

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return mounts, fmt.Errorf("GetBackendType() failed -  %w", err)
	}
	commonVolumeMounts := append(mounts, v1.VolumeMount{
		Name:      "db-cred",
		MountPath: "/srv/k8s/db-cred/"})

	switch backend {
	case "generic":
		return commonVolumeMounts, nil
	case "google":
		return append(commonVolumeMounts,
			v1.VolumeMount{
				Name:      "gcloud-secret",
				MountPath: "/srv/gcloud/"}), nil
	case "amazon":
		return append(commonVolumeMounts,
			v1.VolumeMount{
				Name:      "datastorage-volume",
				MountPath: "/datastorage"}), nil
	default:
		return mounts, errors.New("unable to configure volumeMounts as unknown backend type")
	}
}

func volumes(dbcr *kciv1alpha1.Database) ([]v1.Volume, error) {
	vols := make([]v1.Volume, 0)

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return vols, fmt.Errorf("GetBackendType() failed -  %w", err)
	}
	commonVolumes := append(vols,
		v1.Volume{Name: "db-cred",
			VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: dbcr.Spec.SecretName}}})
	switch backend {
	case "google":
		return append(commonVolumes,
			v1.Volume{Name: "gcloud-secret",
				VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: "google-cloud-storage-bucket-cred"}}}), nil
	case "generic":
		return commonVolumes, nil
	case "amazon":
		return append(commonVolumes,
			v1.Volume{Name: "datastorage-volume",
				VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: ObjectMetaBuilder(dbcr, "volume").Name}}}), nil
	}
	return vols, fmt.Errorf("unsupport backend type: %s", backend)
}

func postgresEnvVars(dbcr *kciv1alpha1.Database) ([]v1.EnvVar, error) {
	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		logrus.Errorf("can not build backup environment variables - %w", err)
		return nil, err
	}

	host, err := getBackupHost(dbcr)
	if err != nil {
		return []v1.EnvVar{}, fmt.Errorf("can not build postgres backup job environment variables - %w", err)
	}

	port := instance.Status.Info["DB_PORT"]

	return []v1.EnvVar{
		{
			Name: "DB_HOST", Value: host,
		},
		{
			Name: "DB_PORT", Value: port,
		},
		{
			Name: "DB_NAME", ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: dbcr.Spec.SecretName},
					Key:                  "POSTGRES_DB",
				},
			},
		},
		{
			Name: "DB_PASSWORD_FILE", Value: "/srv/k8s/db-cred/POSTGRES_PASSWORD",
		},
		{
			Name: "DB_USERNAME_FILE", Value: "/srv/k8s/db-cred/POSTGRES_USER",
		},
		{
			Name: "GCS_BUCKET", Value: instance.Spec.Backup.Bucket,
		},
	}, nil
}

func mysqlEnvVars(dbcr *kciv1alpha1.Database) ([]v1.EnvVar, error) {
	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		logrus.Errorf("can not build backup environment variables - %s", err)
		return nil, err
	}

	host, err := getBackupHost(dbcr)
	if err != nil {
		return []v1.EnvVar{}, fmt.Errorf("can not build mysql backup job environment variables - %s", err)
	}
	port := instance.Status.Info["DB_PORT"]

	return []v1.EnvVar{
		{
			Name: "DB_HOST", Value: host,
		},
		{
			Name: "DB_PORT", Value: port,
		},
		{
			Name: "DB_NAME", ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: dbcr.Spec.SecretName},
					Key:                  "DB",
				},
			},
		},
		{
			Name: "DB_USER", ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: dbcr.Spec.SecretName},
					Key:                  "USER",
				},
			},
		},
		{
			Name: "DB_PASSWORD_FILE", Value: "/srv/k8s/db-cred/PASSWORD",
		},
		{
			Name: "GCS_BUCKET", Value: instance.Spec.Backup.Bucket,
		},
	}, nil
}

func getBackupHost(dbcr *kciv1alpha1.Database) (string, error) {
	var host = ""

	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		return host, err
	}

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return host, err
	}

	switch backend {
	case "google":
		host = "db-" + dbcr.Name + "-svc" //cloud proxy service name
		return host, nil
	case "generic":
		if instance.Spec.Generic.BackupHost != "" {
			return instance.Spec.Generic.BackupHost, nil
		}
		return instance.Spec.Generic.Host, nil
	case "amazon":
		if instance.Spec.Amazon.Generic.BackupHost != "" {
			return instance.Spec.Amazon.Generic.BackupHost, nil
		}
		return instance.Spec.Amazon.Generic.Host, nil
	default:
		return host, errors.New("unknown backend type")
	}
}

func getServiceAccountName(dbcr *kciv1alpha1.Database) (string, error) {
	var account = ""

	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		return account, err
	}

	backend, err := dbcr.GetBackendType()
	if err != nil {
		return account, err
	}

	switch backend {
	case "generic":
		return account, nil
	case "google":
		return account, nil
	case "amazon":
		if instance.Spec.Amazon.ServiceAccountName != "" {
			return instance.Spec.Amazon.ServiceAccountName, nil
		}
		if conf.Instances.Amazon.ServiceAccountName != "" {
			return conf.Instances.Amazon.ServiceAccountName, nil
		}
		return account, nil
	default:
		return account, errors.New("unknown backend type")
	}
}

func getSecurityContext(dbcr *kciv1alpha1.Database) (*v1.PodSecurityContext, error) {
	backend, err := dbcr.GetBackendType()
	if err != nil {
		return nil, err
	}

	switch backend {
	case "generic":
		return nil, nil
	case "google":
		return nil, nil
	case "amazon":
		return getAmazonInstanceSecurityContext(dbcr)
	default:
		return nil, errors.New("unknown backend type")
	}
}

func getAmazonInstanceSecurityContext(dbcr *kciv1alpha1.Database) (*v1.PodSecurityContext, error) {

	instance, err := dbcr.GetInstanceRef()
	if err != nil {
		return nil, err
	}

	instanceFSGroup := (int64)(instance.Spec.Amazon.FSGroup)
	instanceSecurityContext := v1.PodSecurityContext{FSGroup: &instanceFSGroup}
	if instanceFSGroup != -1 {
		return &instanceSecurityContext, nil
	}

	configFSGroup := (int64)(conf.Instances.Amazon.FSGroup)
	configSecurityContext := v1.PodSecurityContext{FSGroup: &configFSGroup}
	if configFSGroup != -1 {
		return &configSecurityContext, nil
	}

	return nil, nil
}

// Contruct a ObjectMeta with information from Database
func ObjectMetaBuilder(db *kciv1alpha1.Database, suffix string) *metav1.ObjectMeta {
	volumeName := db.Namespace + "-" + db.Name + "-" + suffix
	return &metav1.ObjectMeta{
		Name:            volumeName,
		Namespace:       db.Namespace,
		UID:             types.UID(volumeName),
		ResourceVersion: "1",
	}
}
