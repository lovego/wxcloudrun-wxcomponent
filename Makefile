


build:
	docker build  -t three/open-wx:latest .


linuxx86:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t three/open-wx:1.0  .



