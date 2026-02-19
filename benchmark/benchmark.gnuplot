# gnuplot <benchmark.gnuplot

reset

# for windows
set encoding utf8

set terminal pngcairo font "simsun,12" size 1200,850 noenhanced

set style data linespoints
set pointsize 0.8

set output "benchmark.png"
set title "strconv vs byteconv (go1.26)"
set xlabel ""
set ylabel "ns/op"
set xtics rotate by -90

set key right top
set key spacing 1.2
set grid ytics

plot \
	"benchmark_result_byteconv.txt" using 3:xticlabels(1) title "byteconv" with linespoints, \
	"benchmark_result_strconv.txt" using 3:xticlabels(1) title "strconv" with linespoints
