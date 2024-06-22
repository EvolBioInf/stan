export LC_ALL=POSIX
./stan -s 3 -o -r 1501-2000,3501-4000
cat targets/* neighbors/*
