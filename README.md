# Simulate Targets and Neighbors, `stan`
## Description
Simulate sets of target and neighbor DNA sequences under a coalescent
model with defined deletions in the neighbors. These deletions can
then be extracted by programs for neighbor-based marker discovery, for
example [`fur`](https://github.com/evolbioinf/fur).
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Make the Programs
Make sure you've installed the packages `git`, `golang`, `make`, and `noweb`.  
  `$ make`  
  The directory `bin` now contains the binaries, scripts are in
  `scripts`.
## Make the Documentation
Make sure you've installed the packages `git`, `make`, `noweb`, `texlive-science`,
`texlive-pstricks`, `texlive-latex-extra`,
and `texlive-fonts-extra`.  
  `$ make doc`  
  The documentation is now in `doc/stanDoc.pdf`.
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
