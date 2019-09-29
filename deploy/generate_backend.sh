#!/bin/bash 

echo "terraform { 
	backend \"s3\" {
		encrypt = true
		bucket = \"$bucket\"
		region = \"$region\"
		key = \"terraform.state\"
	}
}

variable \"bucket\" { 
	type = \"string\"
	default = \"$bucket\"
}" > backend.tf