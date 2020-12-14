# Database Operator
The original code is [link to github repository ](https://github.com/kloeckner-i/db-operator)

###### Custom Resources Definitions (CRD)
Custom Resources Definitions (CRD) refer to the content of serialized json records held by the k8s api server.

By registering Controllers clients of k8s api can implement the process of reconciliation to desired state by watching for CRD updates..

###### Resources for kci.rocks/v1alpha1
 The operator registers "kci.rocks/v1alpha1" resources with k8s. These are defined in the v1alpha1 package.
[dbinstance_types.go](../pkg/apis/kci/v1alpha1/dbinstance_types.go)
[database_types.go](../pkg/apis/kci/v1alpha1/database_types.go)

###### Controllers for kci.rocks/v1alpha1
 The operator registers controllers for the CRD with k8s,  These are defined in the controller package.
 [database/controller.go](../pkg/controller/database/controller.go)
 [dbinstance/controller.go](../pkg/controller/dbinstance/controller.go)

###### Backup Implementation
The structure of the CRD provided for Database resources, supports the concept of backup. The backup functionality works by creating a k8s cronjob resource.
[backup/cronjob.go](../pkg/controller/database/backup/cronjob.go)
[backup/cronjob_test.go](../pkg/controller/database/backup/cronjob_test.go)

####### Purpose of modifications V1
- Introduce an amazon backend to dbinstance CRD.
- Support use of amazon specific authentication within generated cronjob.
- Support use of volume mount within generated cronjob.
- Support use of persistent volume claim within generated cronjob.
- Support use of fsGroup within generated cronjob.
- Add rclone.conf to docker image
-
####### Purpose of modifications V2
- Support generation of commissioning jobs.
- Introduce an Elasticsearch Engine type to dbInstance.
- Introduce an MongoDB Engine type to dbinstance.
- Introduce an Neo4J Engine type to dbinstance
- Support configuring the name of Elastic backup container.
- Support rolling forward the elasticschema update as currently handled in

###### Extending Backup Implementation
The purpose of the operator is the automation of interaction with amazon hosted datastores, presuming integration with TF, suppport for alternate containers based on the engine type, resolves both comissioning and backup.

There are a number of adaptions required to support Amazon targeted backup as the original implementation was targeted at Google Compute Engine (GCE).
The cronjobs will archive snapshots to a S3 bucket instead of GCE.
The subsequent creation of a point in time restoration should accept a snapshot and apply, giving us a simple method of recovering from data loss.

[pgdump-gcs entrypoint](https://github.com/kloeckner-i/pgdump-gcs/blob/master/entrypoint.sh)

###### AWS  Entrypoint for PSQL backup

```
  1 #!/bin/bash
  2 set -e
  3
  4 echo "Prepare configuration for script"
  5 TIMESTAMP=$(date +%F_%R)
  6 BACKUP_FILE=${DB_NAME}-${TIMESTAMP}.sql.gz
  7 BACKUP_FILE_LATEST=${DB_NAME}-latest.sql.gz
  8 DB_HOST=${DB_HOST:-localhost}
  9 DB_PASSWORD=$(cat ${DB_PASSWORD_FILE})
 10 DB_USER=$(cat ${DB_USERNAME_FILE})
 11
 12 # create login credential file
 13 echo *:5432:*:${DB_USER}:${DB_PASSWORD} >> ~/.pgpass
 14 chmod 0600 ~/.pgpass
 15
 16 pg_dump -F c -Z 9 -h ${DB_HOST} -p 5432 -U ${DB_USER} ${DB_NAME} -f ${BACKUP_FILE}
 17 rclone -v copyto --s3-upload-cutoff 0 --ignore-checksum  ${BACKUP_FILE} s3:${GCS_BUCKET}:${DB_NAME}/${BACKUP_FILE} \
 18 && rclone -v copyto --s3-upload-cutoff 0 --ignore-checksum  ${BACKUP_FILE} s3:${GCS_BUCKET}:${DB_NAME}/${BACKUP_FILE_LATEST}
 19
 20 if test $? -ne 0
 21 then
 22     exit 1;
~
```
