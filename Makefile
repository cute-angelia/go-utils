.PHONY: up

up:
	git pull origin master
	git add .
	git commit -am "update"
	git push origin master
	@echo "\n 发布中..."