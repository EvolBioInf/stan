all: cli.sh
cli.sh: tut.org
	bash ../scripts/org2nw tut.org | notangle -Rcli > cli.sh
test: cli.sh
	bash cli.sh
clean:
	rm -f cli.sh
