ifeq ($(TF_VAR_bucket),)
$(error variable bucket is not set)
endif

ifeq ($(TF_VAR_region),)
$(error variable region is not set)
endif

lambdas = $(wildcard builds/*_lambda)

.PHONY: upload $(lambdas) deploy bucket-up

$(lambdas):
	chmod +x $@
	zip $@.zip $@

bucket-up:
	aws s3api create-bucket --bucket $(TF_VAR_bucket) --region $(TF_VAR_region) --create-bucket-configuration LocationConstraint=$(TF_VAR_region)	

build-dir:
	mkdir -p builds

# each added func/lambda should be <name>_lambda	
build: build-dir
	cd src && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_create_certificate_authority lambdas/cert_create_certificate_authority/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_create_intermediate lambdas/cert_create_intermediate/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_sign_csr lambdas/cert_sign_csr/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/user_get_token lambdas/user_get_token/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/user_authorization lambdas/user_authorization/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/lambda_router lambdas/lambda_router/*.go
	

zip:
	make $(lambdas)

deploy:
	echo  region $(TF_VAR_region)
	cd deploy && \
	bucket=$(TF_VAR_bucket) region=$(TF_VAR_region) ./generate_backend.sh && \
	terraform init && \
	terraform apply
	
	
	
