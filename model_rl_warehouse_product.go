package main

import (
  "sync"
  "encoding/gob"
  "os"
  "unsafe"
  "github.com/golang/glog"
  //"github.com/graphql-go/graphql"
)

// Количество и стоимость товаров на складах

type WarehouseProduct struct {
  Warehouse_ID       string  
  Product_ID         int64  
  Quantity           int
  Cost               int    
}

type WarehouseProductLite struct {
  Quantity           int
  Cost               int    
}

var muPrWH   sync.RWMutex
var memPrWH map[int64]map[string]WarehouseProductLite

var maxx_warehouses int64
var maxx_wproducts int64

func WarehouseProductInit(max_products int64, max_warehouses int64) {
  maxx_warehouses = maxx_warehouses
  maxx_wproducts = max_products
  memPrWH = make(map[int64]map[string]WarehouseProductLite, max_products)
  
  glog.Infof("LOG: Products In Warehouses: sizeof(item) = %d", unsafe.Sizeof(WarehouseProductLite{}))
  glog.Infof("LOG: Products In Warehouses: sizeof(map)  = %d", unsafe.Sizeof(memPrWH))
}

func WarehouseProductAppend(info *WarehouseProduct) {
  muPrWH.Lock()
  if _, ok := memPrWH[info.Product_ID]; !ok {
    memPrWH[info.Product_ID] = make(map[string]WarehouseProductLite, maxx_warehouses)
  }
  var whl WarehouseProductLite
  whl.Quantity = info.Quantity
  whl.Cost = info.Cost
  memPrWH[info.Product_ID][info.Warehouse_ID] = whl
  muPrWH.Unlock()
}

func RlGetProductInWarehouse(product_id int64, warehouse_id string) (*WarehouseProductLite) {
  muPrWH.RLock()
  item, ok := memPrWH[product_id]
  muPrWH.RUnlock()
  if ok {
    res, ok2 := item[warehouse_id]
    if ok2 {
      return &res
    }
  }
  
  return &WarehouseProductLite{Quantity: 0, Cost: 0}
}

func WriteFileQuantityInWarehouse(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memPrWH)
}

func LoadFileQuantityInWarehouse(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memPrWH)
  if err != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename, err)
    return
  }
}
