package main

import (
  "sync"
)

func SaveAll() {
  var wg sync.WaitGroup
  
  wg.Add(1)
  go WriteFileProducts(&wg, "etc/products.gob")
  
  wg.Add(1)
  go WriteFileRegions(&wg, "etc/regions.gob")
  
  wg.Add(1)
  go WriteFileWarehouses(&wg, "etc/warehouses.gob")
  
  wg.Add(1)
  go WriteFileShops(&wg, "etc/shops.gob")
  
  wg.Add(1)
  go WriteFileQuantityInWarehouse(&wg, "etc/cnt_warehouses.gob")
  
  wg.Add(1)
  go WriteFileQuantityInShops(&wg, "etc/cnt_shops.gob")
  
  wg.Wait()
}

func LoadAll() {
  var wg sync.WaitGroup
  
  wg.Add(1)
  go LoadFileProducts(&wg, "etc/products.gob")
  
  wg.Add(1)
  go LoadFileRegions(&wg, "etc/regions.gob")
  
  wg.Add(1)
  go LoadFileWarehouses(&wg, "etc/warehouses.gob")
  
  wg.Add(1)
  go LoadFileShops(&wg, "etc/shops.gob")
  
  wg.Add(1)
  go LoadFileQuantityInWarehouse(&wg, "etc/cnt_warehouses.gob")
  
  wg.Add(1)
  go LoadFileQuantityInShops(&wg, "etc/cnt_shops.gob")
  
  wg.Wait()
}
