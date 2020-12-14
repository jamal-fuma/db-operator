package backup

import (
	"fmt"
	"os"
	"testing"

	kciv1alpha1 "github.com/kloeckner-i/db-operator/pkg/apis/kci/v1alpha1"
	"github.com/kloeckner-i/db-operator/pkg/config"
	v1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
)

func TestGCSBackupCronGsql(t *testing.T) {
	dbcr := &kciv1alpha1.Database{}
	dbcr.Namespace = "TestNS"
	dbcr.Name = "TestDB"
	instance := &kciv1alpha1.DbInstance{}
	instance.Status.Info = map[string]string{"DB_CONN": "TestConnection", "DB_PORT": "1234"}
	instance.Spec.Google = &kciv1alpha1.GoogleInstance{InstanceName: "google-instance-1"}
	dbcr.Status.InstanceRef = instance
	dbcr.Spec.Instance = "staging"
	dbcr.Spec.Backup.Cron = "* * * * *"

	os.Setenv("CONFIG_PATH", "./test/backup_config.yaml")
	conf = config.LoadConfig()

	instance.Spec.Engine = "postgres"
	funcCronObject, err := GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "postgresbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	instance.Spec.Engine = "mysql"
	funcCronObject, err = GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "mysqlbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	assert.Equal(t, "TestNS", funcCronObject.Namespace)
	assert.Equal(t, "TestNS-TestDB-backup", funcCronObject.Name)
	assert.Equal(t, "* * * * *", funcCronObject.Spec.Schedule)
}

func TestGCSBackupCronGeneric(t *testing.T) {
	dbcr := &kciv1alpha1.Database{}
	dbcr.Namespace = "TestNS"
	dbcr.Name = "TestDB"
	instance := &kciv1alpha1.DbInstance{}
	instance.Status.Info = map[string]string{"DB_CONN": "TestConnection", "DB_PORT": "1234"}
	instance.Spec.Generic = &kciv1alpha1.GenericInstance{BackupHost: "slave.test"}
	dbcr.Status.InstanceRef = instance
	dbcr.Spec.Instance = "staging"
	dbcr.Spec.Backup.Cron = "* * * * *"

	os.Setenv("CONFIG_PATH", "./test/backup_config.yaml")
	conf = config.LoadConfig()

	instance.Spec.Engine = "postgres"
	funcCronObject, err := GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
	assert.Equal(t, "postgresbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	instance.Spec.Engine = "mysql"
	funcCronObject, err = GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "mysqlbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	assert.Equal(t, "TestNS", funcCronObject.Namespace)
	assert.Equal(t, "TestNS-TestDB-backup", funcCronObject.Name)
	assert.Equal(t, "* * * * *", funcCronObject.Spec.Schedule)
}

func checkSecurityContextFSGroupEqualValue(fsGroup int64, context *v1.PodSecurityContext, t *testing.T) {
	if -1 == fsGroup {
		checkSecurityContextFSGroupNil(context, t)
	} else {
		checkSecurityContextFSGroupEqual(&fsGroup, context, t)
	}
}

func checkSecurityContextFSGroupEqual(valref *int64, context *v1.PodSecurityContext, t *testing.T) {
	if nil == valref {
		checkSecurityContextFSGroupNil(context, t)
	} else {
		assert.NotNil(t, context, "SecurityContext is nil")
		if nil != context {
			assert.NotNil(t, context.FSGroup, "SecurityContext.FSGroup is Nil")
			if nil != context.FSGroup {
				assert.Equal(t, *valref, *context.FSGroup, "SecurityContext.FSGroup == ")
			}
		}
	}
}

func checkSecurityContextFSGroupNil(context *v1.PodSecurityContext, t *testing.T) {
	assert.Equal(t, (*v1.PodSecurityContext)(nil), context, "SecurityContext should be nil")
	if nil != context {
		assert.Equal(t, (*int64)(nil), context.FSGroup, "SecurityContext.FSGroup should be Nil")
	}
}

func TestGCSBackupCronAmazonServiceAccountFromConfig(t *testing.T) {
	dbcr := &kciv1alpha1.Database{}
	dbcr.Namespace = "TestNS"
	dbcr.Name = "TestDB"
	instance := &kciv1alpha1.DbInstance{}
	instance.Status.Info = map[string]string{"DB_CONN": "TestConnection", "DB_PORT": "1234"}
	instance.Spec.Amazon = &kciv1alpha1.AmazonInstance{Generic: kciv1alpha1.GenericInstance{BackupHost: "slave.test"}, FSGroup: -1}

	dbcr.Status.InstanceRef = instance
	dbcr.Spec.Instance = "staging"
	dbcr.Spec.Backup.Cron = "* * * * *"

	os.Setenv("CONFIG_PATH", "./test/backup_config.yaml")
	conf = config.LoadConfig()

	instance.Spec.Engine = "postgres"
	funcCronObject, err := GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "backup", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
	assert.Equal(t, "postgresbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	securityContext := funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext
	fsGroup := conf.Instances.Amazon.FSGroup
	checkSecurityContextFSGroupEqualValue(fsGroup, securityContext, t)
	assert.Equal(t, fsGroup, *funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext.FSGroup)
	expectedVolumeMounts := []v1.VolumeMount{{Name: "db-cred", MountPath: "/srv/k8s/db-cred/"}, {Name: "datastorage-volume", MountPath: "/datastorage"}}

	actualVolumeMounts := funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts
	assert.Equal(t, expectedVolumeMounts, funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts)
	assert.Equal(t, len(expectedVolumeMounts), len(actualVolumeMounts), "expected that expectedVolumeMounts and actualVolumeMounts would have same length")
	for idx, volMount := range actualVolumeMounts {
		assert.Equal(t, volMount.Name, expectedVolumeMounts[idx].Name)
		assert.Equal(t, volMount.MountPath, expectedVolumeMounts[idx].MountPath)
	}

	podSpec := funcCronObject.Spec.JobTemplate.Spec.Template.Spec
	securityContext = podSpec.SecurityContext
	containers := podSpec.Containers
	assert.Equal(t, "backup", podSpec.ServiceAccountName)
	assert.Equal(t, "postgresbackupimage:latest", containers[0].Image)
	assert.Equal(t, expectedVolumeMounts, containers[0].VolumeMounts)
	checkSecurityContextFSGroupEqualValue(fsGroup, podSpec.SecurityContext, t)
	assert.Equal(t, fsGroup, *podSpec.SecurityContext.FSGroup)

	instance.Spec.Engine = "mysql"
	funcCronObject, err = GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "mysqlbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	securityContext = funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext
	fsGroup = conf.Instances.Amazon.FSGroup
	checkSecurityContextFSGroupEqualValue(fsGroup, securityContext, t)
	assert.Equal(t, fsGroup, *funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext.FSGroup)

	assert.Equal(t, "backup", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
	assert.Equal(t, "TestNS", funcCronObject.Namespace)
	assert.Equal(t, "TestNS-TestDB-backup", funcCronObject.Name)
	assert.Equal(t, "* * * * *", funcCronObject.Spec.Schedule)
}

func TestGCSBackupCronAmazonServiceAccountFromInstance(t *testing.T) {
	dbcr := &kciv1alpha1.Database{}
	dbcr.Namespace = "TestNS"
	dbcr.Name = "TestDB"
	instance := &kciv1alpha1.DbInstance{}
	instance.Status.Info = map[string]string{"DB_CONN": "TestConnection", "DB_PORT": "1234"}
	instance.Spec.Amazon = &kciv1alpha1.AmazonInstance{Generic: kciv1alpha1.GenericInstance{BackupHost: "slave.test"}, ServiceAccountName: "backup01", FSGroup: 456}
	dbcr.Status.InstanceRef = instance
	dbcr.Spec.Instance = "staging"
	dbcr.Spec.Backup.Cron = "* * * * *"

	os.Setenv("CONFIG_PATH", "./test/backup_config.yaml")
	conf = config.LoadConfig()

	instance.Spec.Engine = "postgres"
	funcCronObject, err := GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	securityContext := funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext
	fsGroup := instance.Spec.Amazon.FSGroup
	checkSecurityContextFSGroupEqualValue(fsGroup, securityContext, t)
	assert.Equal(t, fsGroup, *funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext.FSGroup)

	assert.Equal(t, "backup01", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
	assert.Equal(t, "postgresbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	instance.Spec.Engine = "mysql"
	funcCronObject, err = GCSBackupCron(dbcr)
	if err != nil {
		fmt.Print(err)
	}

	assert.Equal(t, "mysqlbackupimage:latest", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)

	securityContext = funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext
	fsGroup = instance.Spec.Amazon.FSGroup
	checkSecurityContextFSGroupEqualValue(fsGroup, securityContext, t)
	assert.Equal(t, fsGroup, *funcCronObject.Spec.JobTemplate.Spec.Template.Spec.SecurityContext.FSGroup)

	expectedVolumeMounts := []v1.VolumeMount{{Name: "db-cred", MountPath: "/srv/k8s/db-cred/"}, {Name: "datastorage-volume", MountPath: "/datastorage"}}
	actualVolumeMounts := funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts
	assert.Equal(t, expectedVolumeMounts, funcCronObject.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts)
	assert.Equal(t, len(expectedVolumeMounts), len(actualVolumeMounts), "expected that expectedVolumeMounts and actualVolumeMounts would have same length")
	for idx, volMount := range actualVolumeMounts {
		assert.Equal(t, volMount.Name, expectedVolumeMounts[idx].Name)
		assert.Equal(t, volMount.MountPath, expectedVolumeMounts[idx].MountPath)
	}

	assert.Equal(t, "backup01", funcCronObject.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName)
	assert.Equal(t, "TestNS", funcCronObject.Namespace)
	assert.Equal(t, "TestNS-TestDB-backup", funcCronObject.Name)
	assert.Equal(t, "* * * * *", funcCronObject.Spec.Schedule)
}
