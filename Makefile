# Nombre de la imagen
IMAGE_NAME = cloudbuilders
TAG = latest

# ---- Go ----
run:
	go run ./cmd/webterm

# ---- Docker ----
docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-run:
	docker run -p 8080:8080 $(IMAGE_NAME):$(TAG)

docker-push:
	docker tag $(IMAGE_NAME):$(TAG) tuusuario/$(IMAGE_NAME):$(TAG)
	docker push tuusuario/$(IMAGE_NAME):$(TAG)

# ---- Kubernetes ----
k8s-apply:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/

k8s-status:
	kubectl get pods,svc
