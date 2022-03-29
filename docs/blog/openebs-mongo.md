# Using this CloudImage

This is a CloudImage.

```
sealos run \
   -e openebs-basedir=/data -e mongo-replicaCount=3 \
   fanux/kubernetes:v1.23.5 fanux/openebs:latest fanux/mongo:latest
```

# Config storage Dir

storage-class.yaml:
```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-hostpath
  annotations:
    openebs.io/cas-type: local
    cas.openebs.io/config: |
      - name: StorageType
        value: hostpath
      - name: BasePath
        value: /var/local-hostpath # Host path storage dir
provisioner: openebs.io/local
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

Edit the BasePath value, then:
```
kubectl apply -f storage-class.yaml
```

# Using storage

> Create PVC local-hostpath-pvc.yaml:

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: local-hostpath-pvc
spec:
  storageClassName: openebs-hostpath
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5G
```

```
kubectl apply -f local-hostpath-pvc.yaml:
```

Now the PVC not bind
```
kubectl get pvc local-hostpath-pvc
NAME                 STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS       AGE
local-hostpath-pvc   Pending                                      openebs-hostpath   3m7s
```

> Create a pod to consume OpenEBS

```
apiVersion: v1
kind: Pod
metadata:
  name: hello-local-hostpath-pod
spec:
  volumes:
  - name: local-storage
    persistentVolumeClaim:
      claimName: local-hostpath-pvc
  containers:
  - name: hello-container
    image: busybox
    command:
       - sh
       - -c
       - 'while true; do echo "`date` [`hostname`] Hello from OpenEBS Local PV." >> /mnt/store/greet.txt; sleep $(($RANDOM % 5 + 300)); done'
    volumeMounts:
    - mountPath: /mnt/store
      name: local-storage
```

```
kubectl apply -f local-hostpath-pod.yaml
kubectl exec hello-local-hostpath-pod -- cat /mnt/store/greet.txt
```

```
kubectl get pvc local-hostpath-pvc
NAME                 STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS       AGE
local-hostpath-pvc   Bound    pvc-864a5ac8-dd3f-416b-9f4b-ffd7d285b425   5G         RWO            openebs-hostpath   28m
```

Checkout the bound pvc

```
kubectl get pv pvc-864a5ac8-dd3f-416b-9f4b-ffd7d285b425 -o yaml
```

# Clean up

```
kubectl delete pod hello-local-hostpath-pod
kubectl delete pvc local-hostpath-pvc
kubectl delete sc local-hostpath
```

```
kubectl get pv
```

[link](https://openebs.io/docs/user-guides/localpv-hostpath)

# Trace your data

1. Get the pod pvc name:

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ openebs]# kubectl get pod hello-local-hostpath-pod-4 -oyaml|grep claimName
      claimName: local-hostpath-pvc-4
```

2. Get the pv nodename and pvname

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ openebs]# kubectl get pvc local-hostpath-pvc-4 -oyaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
...
    volume.kubernetes.io/selected-node: iz2ze0qiwmjj4p5rncuhhoz
...
  name: local-hostpath-pvc-4
...
  storageClassName: local-hostpath
  volumeName: pvc-056c7781-c9b3-46f6-aa6e-a3a2d72456d6
...
```

We got the data selected-node is: iz2ze0qiwmjj4p5rncuhhoz

You use the storageClass is: local-hostpath

Show the storageClass detail:

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ openebs]# kubectl get sc local-hostpath -oyaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    cas.openebs.io/config: |
      - name: StorageType
        value: hostpath
      - name: BasePath
        value: /data
    openebs.io/cas-type: local
...
provisioner: openebs.io/local
reclaimPolicy: Delete
```

So the data dir is `/data`.

So your data is storage in: iz2ze0qiwmjj4p5rncuhhoz:/data/pvc-056c7781-c9b3-46f6-aa6e-a3a2d72456d6

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ openebs]# kubectl get node -owide
NAME                      STATUS   ROLES                  AGE   VERSION   INTERNAL-IP     EXTERNAL-IP   OS-IMAGE                KERNEL-VERSION              CONTAINER-RUNTIME
iz2ze0qiwmjj4p5rncuhhoz   Ready    <none>                 29h   v1.22.0   172.17.83.145   <none>        CentOS Linux 7 (Core)   3.10.0-693.2.2.el7.x86_64   containerd://1.4.3
```

```
ssh root@172.17.83.145
[root@iZ2ze0qiwmjj4p5rncuhhoZ pvc-056c7781-c9b3-46f6-aa6e-a3a2d72456d6]# cd /data/pvc-056c7781-c9b3-46f6-aa6e-a3a2d72456d6 && ls
greet.txt
```

# Using openebs as mongo storage class

```
git clone https://github.com/bitnami/charts
```

Config values (replicaset mod):
```
architecture=replicaset
replicaCount=3
externalAccess.enabled=true
externalAccess.service.type=NodePort
externalAccess.service.nodePorts[0]='31001'
externalAccess.service.nodePorts[1]='31002'
externalAccess.service.nodePorts[1]='31003'
```

Change StorageClass:

```
storageClass: "local-hostpath"
```

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ mongodb]# cd bitnami/mongodb && helm install mongo-test .
NAME: mongo-test
LAST DEPLOYED: Tue Mar 29 16:18:08 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mongodb
CHART VERSION: 11.1.3
APP VERSION: 4.4.13

** Please be patient while the chart is being deployed **

MongoDB&reg; can be accessed on the following DNS name(s) and ports from within your cluster:

    mongo-test-mongodb-0.mongo-test-mongodb-headless.default.svc.cluster.local:27017
    mongo-test-mongodb-1.mongo-test-mongodb-headless.default.svc.cluster.local:27017
    mongo-test-mongodb-2.mongo-test-mongodb-headless.default.svc.cluster.local:27017
    mongo-test-mongodb-3.mongo-test-mongodb-headless.default.svc.cluster.local:27017

To get the root password run:

    export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace default mongo-test-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

To connect to your database, create a MongoDB&reg; client container:

    kubectl run --namespace default mongo-test-mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:4.4.13-debian-10-r25 --command -- bash

Then, run the following command:
    mongo admin --host "mongo-test-mongodb-0.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-1.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-2.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-3.mongo-test-mongodb-headless.default.svc.cluster.local:27017" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

To connect to your database nodes from outside, you need to add both primary and secondary nodes hostnames/IPs to your Mongo client. To obtain them, follow the instructions below:

    MongoDB&reg; nodes domain: you can reach MongoDB&reg; nodes on any of the K8s nodes external IPs.

        kubectl get nodes -o wide

    MongoDB&reg; nodes port: You will have a different node port for each MongoDB&reg; node. You can get the list of configured node ports using the command below:

        echo "$(kubectl get svc --namespace default -l "app.kubernetes.io/name=mongodb,app.kubernetes.io/instance=mongo-test,app.kubernetes.io/component=mongodb,pod" -o jsonpath='{.items[*].spec.ports[0].nodePort}' | tr ' ' '\n')"
```

Check your pod:

```
[root@iZ2ze0qiwmjj4p5rncuhhrZ mongodb]# kubectl get pod
NAME                           READY   STATUS      RESTARTS      AGE
mongo-test-mongodb-0           1/1     Running     0             49m
mongo-test-mongodb-1           1/1     Running     0             49m
mongo-test-mongodb-2           0/1     Running     1 (90s ago)   48m
mongo-test-mongodb-arbiter-0   1/1     Running     0             49m
```

Check the pvc:
```
[root@iZ2ze0qiwmjj4p5rncuhhrZ mongodb]# kubectl get pvc
NAME                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS       AGE
datadir-mongo-test-mongodb-0   Bound    pvc-5bddcedc-eb0c-41ed-a230-f7c953bc537f   8Gi        RWO            local-hostpath     52m
datadir-mongo-test-mongodb-1   Bound    pvc-c187a64a-c3e6-4e4b-9669-c01e30af1dc7   8Gi        RWO            local-hostpath     51m
datadir-mongo-test-mongodb-2   Bound    pvc-b845673f-2297-40ed-b013-603b09acbdd2   8Gi        RWO            local-hostpath     51m
```

Using Client pod to test your cluster:
```
kubectl run --namespace default mongo-test-mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:4.4.13-debian-10-r25 --command -- bash
```

Run mongo cli:
```
 mongo admin --host "mongo-test-mongodb-0.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-1.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-2.mongo-test-mongodb-headless.default.svc.cluster.local:27017,mongo-test-mongodb-3.mongo-test-mongodb-headless.default.svc.cluster.local:27017" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

Implicit session: session { "id" : UUID("25ae50c1-932f-416d-b164-871c9144118d") }
MongoDB server version: 4.4.13
---
The server generated these startup warnings when booting: 
        2022-03-29T08:18:28.221+00:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
        2022-03-29T08:18:28.460+00:00: /sys/kernel/mm/transparent_hugepage/enabled is 'always'. We suggest setting it to 'never'
        2022-03-29T08:18:28.460+00:00: /sys/kernel/mm/transparent_hugepage/defrag is 'always'. We suggest setting it to 'never'
---
---
        Enable MongoDB's free cloud-based monitoring service, which will then receive and display
        metrics about your deployment (disk utilization, CPU, operation statistics, etc).

        The monitoring data will be available on a MongoDB website with a unique URL accessible to you
        and anyone you share the URL with. MongoDB may use this information to make product
        improvements and to suggest MongoDB products and deployment options to you.

        To enable free monitoring, run the following command: db.enableFreeMonitoring()
        To permanently disable this reminder, run the following command: db.disableFreeMonitoring()
---
rs0:PRIMARY> 
```

```
rs0:PRIMARY> help
	db.help()                    help on db methods
	db.mycoll.help()             help on collection methods
	sh.help()                    sharding helpers
	rs.help()                    replica set helpers
	help admin                   administrative help
	help connect                 connecting to a db help
	help keys                    key shortcuts
	help misc                    misc things to know
	help mr                      mapreduce

	show dbs                     show database names
	show collections             show collections in current database
	show users                   show users in current database
	show profile                 show most recent system.profile entries with time >= 1ms
	show logs                    show the accessible logger names
	show log [name]              prints out the last segment of log in memory, 'global' is default
	use <db_name>                set current database
	db.mycoll.find()             list objects in collection mycoll
	db.mycoll.find( { a : 1 } )  list objects in mycoll where a == 1
	it                           result of the last line evaluated; use to further iterate
	DBQuery.shellBatchSize = x   set default number of items to display on shell
	exit                         quit the mongo shell
rs0:PRIMARY> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
```
