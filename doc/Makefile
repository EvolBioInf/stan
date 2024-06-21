SRC1 =   rad.tex stan.tex util.tex tut.tex
SRC2 =   header.tex intro.tex stanDoc.tex

date = $(shell bash ../scripts/getDate.sh)
version = $(shell bash ../scripts/getVersion.sh)

all: stanDoc.pdf
stanDoc.pdf: stanDoc.tex $(SRC1) $(SRC2)
	echo $(date) | tr '_' ' ' > date.txt
	echo $(version) | tr '-' ' ' | awk '{printf "%s", $$1; if ($$2) printf "-%s", $$2; printf "\n"}' > version.txt
	latex stanDoc.tex
	bibtex stanDoc
	latex stanDoc.tex
	latex stanDoc.tex
	dvips stanDoc -o -q
	ps2pdf -dALLOWPSTRANSPARENCY stanDoc.ps
publish: stanDoc.pdf
	if [ -d /mnt/c/Users/haubold/ownCloud\ -\ haubold@evolbio.mpg.de@owncloud.gwdg.de/docs/ ]; then \
		cp stanDoc.pdf /mnt/c/Users/haubold/ownCloud\ -\ haubold@evolbio.mpg.de@owncloud.gwdg.de/docs/; \
	fi
rad.tex: ../rad/rad.org
	bash ../scripts/org2nw ../rad/rad.org | noweave -n -x | sed 's/_test/\\_test/g' > rad.tex
stan.tex: ../stan/stan.org 
	bash ../scripts/org2nw ../stan/stan.org | noweave -n -x | sed 's/_test/\\_test/g' > stan.tex
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/g' > util.tex
tut.tex: ../tut/tut.org
	bash ../scripts/org2nw ../tut/tut.org | noweave -n -x | sed 's/_test/\\_test/g' > tut.tex
clean:
	rm -f stanDoc.pdf stanDoc.ps stanDoc.dvi stanDoc.aux $(SRC1) 
