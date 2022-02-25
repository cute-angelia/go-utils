.PHONY: up tag

up:
	git pull origin master
	git add .
	git commit -am "update"
	git push origin master
	@echo "\n 发布中..."

tag:
	git pull origin master
	git add .
	git commit -am "获得一个干净的不带参数的地址"
	git push origin master
	git tag v1.6.91
	git push --tags
	@echo "\n tags 发布中..."
