#!/bin/zsh

timeout=1800

#go build; 
results=results2
rm -fr $results
mkdir -p $results

for amo in 0 1 
do
	for same in 0 1 
	do
		for x in $1/*.pb; 
		do 
		echo ./bule -d -cat 2 -rewrite-same=$same -solve -amo-chain=$amo -amo-reuse=1 -timeout $timeout -f $x" > "$results/$(basename $x .pb)-amo-$amo-same-$same.log
		done 
	done 
done


#for x in $1/*.pb; 
#do 
#    echo ./bule -cat 1 -opt-bound $bound -solve -timeout $timeout -solver=minisat -f $x" | grep 'xxx:' >> "$results/$(basename $x .pb)-cat-1-opt-0-amo-0.log
#done 
#
#for amo in 0 1
#do 
#    for opt in 0 1 
#    do 
#    	for same in 0 1 
#    	do 
#        	for x in $1/*.pb; 
#        	do 
#		echo ./bule -cat 2 -opt-bound $bound -rewrite-same=$same -solve -amo-reuse=$amo -opt-rewrite=$opt -timeout $timeout -solver=minisat -f $x" | grep 'xxx:' >> "$results/$(basename $x .pb)-cat-2-opt-$opt-amo-$amo-same-$same-bound-$bound.log
#        	done 
#    	done 
#    done 
#done
#
