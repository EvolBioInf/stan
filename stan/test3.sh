./stan -s 3 -o -N otherNeighbors
cat targets/* otherNeighbors/*
rm -r otherNeighbors
