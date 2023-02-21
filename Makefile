.PHONY: up tag

up:
	git pull origin v1
	git add .
	git commit -am "update"
	git push origin v1
	@echo "\n 发布中..."

tag:
	git pull origin master
	git add .
	git commit -am "update"
	git push origin master
	git tag v1.7.72
	git push --tags
	@echo "\n tags 发布中..."