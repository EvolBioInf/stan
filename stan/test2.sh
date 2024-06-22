export LC_ALL=POSIX
./stan -s 3 -o -M 0.01
cat targets/* neighbors/*
