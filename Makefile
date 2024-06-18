


build:
	docker build  -t lchjczw/open-wx:latest .


linuxx86:
	docker buildx build --platform linux/amd64 -t lchjczw/open-wx:latest  .



push:
	docker push lchjczw/open-wx:latest


release:linuxx86 push