./stan -s 3 -c > r1.txt
for a in $(seq 1 9); do
    i=$(($a+1))
    bash test$a.sh > r$i.txt
done
