set terminal png

set output '10.12-max.png'
set title 'Maximum hyperkit memory usage reported in Activity Monitor on 10.12'
set xlabel 'VM size in settings / GiB'
set ylabel 'Hyperkit memory / GiB'
plot 'graph.dat' using 1:2 with linespoints title 'Memory', \
     'graph.dat' using 1:3 with linespoint title 'Real memory'
