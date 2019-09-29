ifeq ($(bucket),)
$(error variable bucket is not set)
endif

ifeq ($(region),)
$(error variable region is not set)
endif

lambdas = $(wildcard builds/*_lambda)

.PHONY: upload $(lambdas) deploy

$(lambdas):
	chmod +x $@
	zip $@.zip $@

bucket-up:
	aws s3api create-bucket --bucket $(bucket) --region $(region) --create-bucket-configuration LocationConstraint=$(region)	

build-dir:
	mkdir -p builds

# each added func/lambda should be <name>_lambda	
build: build-dir
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o builds/create_ca_lambda lambdas/create_ca/*.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o builds/create_intermediate_lambda lambdas/create_intermediate/*.go
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o builds/sign_user_certificate_lambda lambdas/sign_user_certificate/*.go
	#go build -o builds/get_ca_intermediate_lambda lambdas/get_ca_intermediate/*.go

zip:
	make $(lambdas)

upload:
	make $(lambdas) && \
	aws s3 cp builds/ s3://$(bucket)/ --recursive --exclude "*" --include "*.zip"

deploy:
	cd deploy && \
	./generate_backend.sh bucket=$(bucket) region=$(region) && \
	terraform init && \
	terraform apply
	
	
	