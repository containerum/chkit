
### create deployment

**Aliases**   :
  depl, deployments, deploy
**Usage**     :
 Create a new deployment.
Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.
**Example**   :
  
**Flags**     :
  + configmap  : container configmap, CONTAINER_NAME@CONFIGMAP_NAME@MOUNTPATH in case of multiple containers or
  CONFIGMAP_NAME@MOUNTPATH or CONFIGMAP_NAME in case of one container.
  If MOUNTPATH is omitted, then use /etc/CONFIGMAP_NAME as mountpath
  + cpu  : container memory limit, mCPU,
  CONTAINER_NAME@CPU in case of multiple containers or CPU in case of one container
  + env  : container environment variable,
  CONTAINER_NAME@KEY:VALUE in case of multiple containers or KEY:VALUE in case of one container
  + file  : file with configmap data, .json, .yaml, .yml, optional
  + force f : suppress confirmation, optional
  + image  : container image,
  CONTAINER_NAME@IMAGE in case of multiple containers or IMAGE in case of one container
  + memory  : container memory limit, Mb,
  CONTAINER_NAME@MEMORY in case of multiple containers or MEMORY in case of one container
  + name  : deployment name, optional
  + replicas  : deployment replicas, optional
  + volume  : container volume,
  CONTAINER_NAME@VOLUME_NAME@MOUNTPATH in case of multiple containers or
  VOLUME_NAME@MOUNTPATH or VOLUME_NAME in case of one container.
  If MOUNTPATH is omitted, then use /mnt/VOLUME_NAME as mountpath
  
**Subcommand**:
  
