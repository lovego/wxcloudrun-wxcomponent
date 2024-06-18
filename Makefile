


build:
	docker build  -t lchjczw/open-wx:latest .


linuxx86:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t lchjczw/open-wx:1.0  .



push:
	docker push lchjczw/open-wx:latest