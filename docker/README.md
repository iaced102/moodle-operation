```
cd docker
docker build . -t moodle
docker tag $IMAGE_ID cr-hn-1.bizflycloud.vn/c579f3ab86094f2c97f36b0877ead59a/moodle:custom
docker push cr-hn-1.bizflycloud.vn/c579f3ab86094f2c97f36b0877ead59a/moodle:custom
cd ../deploy
kubectl apply -f . -n $NS
```
