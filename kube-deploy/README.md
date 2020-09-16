# Kubernetes deployment instructions

This project requires helm charts for kafka and elassandra. Elassandra integrates Elasticsearch into Apache Cassandra.

## Helm charts installation

- Install kafka: `helm install <kafka-release-name> bitnami/kafka`
- Install elassandra: `helm install my-release-elassandra --set image.tag="6.8.4.10" strapdata/elassandra`

## Services deployments

The project has 3 deployments:

- The reddit streaming service `reddit-kafka`
- The cassandra storage service ` reddit-storage`
- The classification service `reddit-classifier`


for each of the deployment you need to push a docker image. In my case i'll use google cloud:

- `gcloud builds submit --tag gcr.io/<project-id>/reddit-kafka`
- `gcloud builds submit --tag gcr.io/<project-id>/reddit-storage`
- `gcloud builds submit --tag gcr.io/<project-id>/reddit-classifier`

you can than create the deployments:

- `kubectl apply -f reddit-kafka-deployment.yaml`
- `kubectl apply -f reddit-storage-deployment.yaml`
- `kubectl apply -f reddit-classifier-deployment.yaml`

