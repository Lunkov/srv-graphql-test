package main

import (
  "sync"
  "encoding/gob"
  "os"
  "unsafe"
  "github.com/golang/glog"
  //"github.com/graphql-go/graphql"
)

// Количество и стоимость товаров в магазинах

type ShopProduct struct {
  Shop_ID       string  
  Product_ID    int64  
  Quantity      int
  Cost          int    
}

type ShopProductLite struct {
  Quantity      int
  Cost          int    
}

var muPrSP   sync.RWMutex

var memPrSP = make(map[int64]map[string]ShopProductLite, 100000)
var maxx_shops int64
var maxx_products int64

func ShopProductInit(max_products int64, max_shops int64) {
  maxx_shops = max_shops
  maxx_products = max_products
  memPrSP = make(map[int64]map[string]ShopProductLite, max_products)
  
  glog.Infof("LOG: Products In Shops: sizeof(item) = %d", unsafe.Sizeof(ShopProductLite{}))
  glog.Infof("LOG: Products In Shops: sizeof(map)  = %d", unsafe.Sizeof(memPrSP))
}

func ShopProductAppend(info *ShopProduct) {
  muPrSP.Lock()
  if _, ok := memPrSP[info.Product_ID]; !ok {
    memPrSP[info.Product_ID] = make(map[string]ShopProductLite, maxx_products)
  }
  var sgl ShopProductLite
  sgl.Quantity = info.Quantity
  sgl.Cost = info.Cost  
  memPrSP[info.Product_ID][info.Shop_ID] = sgl
  muPrSP.Unlock()
}

func RlGetProductInShop(product_id int64, shop_id string) (*ShopProductLite) {
  muPrSP.RLock()
  item, ok := memPrSP[product_id]
  muPrSP.RUnlock()
  if ok {
    res, ok2 := item[shop_id]
    if ok2 {
      return &res
    }
  }
  
  return &ShopProductLite{Quantity: 0, Cost: 0}
}

func WriteFileQuantityInShops(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memPrSP)
}

func LoadFileQuantityInShops(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memPrSP)
  if err != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename, err)
    return
  }
}
