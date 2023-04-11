build:
	sudo docker build --tag cbr-currency-app .
run:
	sudo docker run cbr-currency-app

.PHONY: build, run