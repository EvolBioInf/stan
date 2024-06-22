export LC_ALL=POSIX
./stan -s 3 -o -T otherTargets
cat otherTargets/* neighbors/*
rm -r otherTargets
