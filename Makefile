
docker-build:
	docker build  --no-cache -t gcr.io/septapig/goetc:test -f Dockerfile .

push:
	docker push gcr.io/septapig/goetc:test

build:
	go build -v .

run:
	docker run --name goetc --rm -it -p 3000:3000  gcr.io/septapig/goetc:test


deploy:
	gcloud beta run deploy goetc  --image gcr.io/septapig/goetc:test --platform managed \
            --allow-unauthenticated --project septapig \
            --vpc-connector=cloudvpc-east \
            --vpc-egress=all \
            --region us-east1 --port 3000 --max-instances 3  --memory 124Mi


