run:
	./scripts/build_and_run_no_docker.sh

run-docker:
	docker build -t spannerprober .
	docker run --rm spannerprober -i -t --name spannerprober

build-upload-image:
	./scripts/build_and_upload_docker_image_with_tag.sh

create-table:
	./scripts/create_table.sh

fmt:
	./scripts/fmt.sh
