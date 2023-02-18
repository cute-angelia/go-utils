.PHONY: up tag

up:
	git pull origin v1
	git add .
	git commit -am "update"
	git push origin v1
	@echo "\n 发布中..."