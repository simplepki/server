ifeq ($(bucket),)
$(error variable bucket is not set)
endif

ifeq ($(region),)
$(error variable region is not set)
endif

lambdas = $(wildcard builds/*_lambda)

.PHONY: upload $(lambdas) deploy bucket-up

$(lambdas):
	chmod +x $@
	zip $@.zip $@

bucket-up:
	aws s3api create-bucket --bucket $(bucket) --region $(region) --create-bucket-configuration LocationConstraint=$(region)	

build-dir:
	mkdir -p builds

# each added func/lambda should be <name>_lambda	
build: build-dir
	cd src && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_create_certificate_authority lambdas/cert_create_certificate_authority/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_create_intermediate lambdas/cert_create_intermediate/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/cert_sign_csr lambdas/cert_sign_csr/*.go && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ../builds/user_sign_in lambdas/user_sign_in/*.go
	#GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o builds/cert_get_certificate_chain lambdas/cert_get_certificate_chain/*.go
	

zip:
	make $(lambdas)

deploy:
	cd deploy && \
	./generate_backend.sh bucket=$(bucket) region=$(region) && \
	terraform init && \
	terraform apply
	
	
	