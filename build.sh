version=$1
dockerRegistry=${2:-pulkitdocker1234}

echo "============== Building image rate-limiter:$version =============="
docker build -t rate-limiter:$version .

echo "============== Tagging image $dockerRegistry/rate-limiter:$version =============="
docker tag rate-limiter:$version $dockerRegistry/rate-limiter:$version

echo "============== Pushing image to registry: $dockerRegistry/rate-limiter:$version =============="
docker push $dockerRegistry/rate-limiter:$version
err=$?
if [ $err -ne 0 ]; then
  echo "Failed to push image to registry"
  exit 1
fi

echo "============== Removing local images =============="
docker rmi $dockerRegistry/rate-limiter:$version
docker rmi rate-limiter:$version

echo "============== Pushed $dockerRegistry/rate-limiter:$version to registry =============="