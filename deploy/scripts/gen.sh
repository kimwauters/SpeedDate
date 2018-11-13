# Generate user-service
cd $GOPATH/src/github.com/proepkes/speeddate/usersvc
goa gen github.com/proepkes/speeddate/usersvc/design
goa example github.com/proepkes/speeddate/usersvc/design

# Generate auth-service
cd $GOPATH/src/github.com/proepkes/speeddate/authsvc
goa gen github.com/proepkes/speeddate/authsvc/design
goa example github.com/proepkes/speeddate/authsvc/design

# Build & Tag & push userservice
docker build --tag=proepkes/usersvc:dev -f $GOPATH/src/github.com/proepkes/speeddate/deploy/docker/Dockerfile.usersvc $GOPATH/src/github.com/proepkes/speeddate/usersvc
docker push proepkes/usersvc:dev

# Generate k8s-yaml
helm template deploy/charts/ -f deploy/charts/values.yaml > deploy/k8s/speeddate.yaml
kubectl apply -f deploy/k8s/speeddate.yaml
