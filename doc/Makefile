SRC1 =   stan.tex util.tex tut.tex
SRC2 =   header.tex intro.tex stanDoc.tex

all: stanDoc.pdf
stanDoc.pdf: stanDoc.tex $(SRC1) $(SRC2)
	latex stanDoc.tex
	bibtex stanDoc
	latex stanDoc.tex
	latex stanDoc.tex
	dvips stanDoc -o -q
	ps2pdf -dALLOWPSTRANSPARENCY stanDoc.ps

stan.tex: ../stan/stan.org 
	bash ../scripts/org2nw ../stan/stan.org | noweave -n -x | sed 's/_test/\\_test/g' > stan.tex
util.tex: ../util/util.org
	bash ../scripts/org2nw ../util/util.org | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/g' > util.tex
tut.tex: ../tut/tut.org
	bash ../scripts/org2nw ../tut/tut.org | noweave -n -x | sed 's/_test/\\_test/g' > tut.tex

clean:
	rm -f stanDoc.pdf stanDoc.ps stanDoc.dvi stanDoc.aux $(SRC1) 
