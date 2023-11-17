.PHONY: deploy start help

help:
	@echo "Available targets:"
	@echo "  deploy      builds and deploys to AWS"
	@echo "  start       starts npm server"

deploy:
	aws s3 sync ./ s3://jz.io/ --acl public-read
	aws cloudfront create-invalidation --distribution-id E9KIROWWL1OJ2 --paths "/*"

start:
	npm start

