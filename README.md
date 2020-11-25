# Before tests
```
sudo apt-get install libcanberra-gtk-module
sudo apt-get install libcanberra-gtk-module libcanberra-gtk3-module
sudo apt-get install graphviz
```

# Tests
About tools: https://blog.golang.org/pprof
```
go test -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.out -memprofile=mem.out
```

```
go tool pprof ./mem.out
go tool pprof ./cpu.out
```

```
top10
web mallocgc
```
# Tests results

Intel® Core™ i5-4210U CPU @ 1.70GHz × 4

* cRegions =  1000 // Количество регионов
* cWH   =    10000 // Количество складов
* cSP   =     5000 // Количество магазинов
* cPR   =   200000 // Количество товаров
* cPRWH =  5000000 // Количество товаров на складах
* cPRSP =  1000000 // Количество товаров в магазинах
