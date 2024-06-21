# [`stan`](https://owncloud.gwdg.de/index.php/s/85gRtxIFrkYZiuj)
## Description
Simulate sets of target and neighbor DNA sequences under a coalescent
model with defined markers in the targets. By default, the markers are
deleted in the neighbors, but alternatively, they can be mutated with
a higher rate than the background in the targets. Markers can then be
extracted by programs for neighbor-based marker discovery, for example
[`fur`](https://github.com/evolbioinf/fur).

## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`

## Make the Programs
Setup the environment and construct the binaries.

```
bash scripts/setup.sh
make
```

The directory `bin` now contains the binaries.

## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
