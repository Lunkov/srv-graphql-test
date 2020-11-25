package main

import (
)

type Product4Sale struct {
  Prod         Product        `json:"product"`
  Quantity     int            `json:"quantity"`
  Shops        []Shop         `json:"shops"`
  Warehouses   []Warehouse    `json:"warehouses"`
}

// Основная функция: получить остатки по товару в зависимости от региона
func GetProduct4Sale(product_id int64, region_id int64) Product4Sale {
  var res Product4Sale
  res.Quantity = 0
  p := GetProductByID(product_id)
  if p == nil {
    return res
  }
  res.Prod = *p
  for _, shop_code := range GetShopByRegionID(region_id) {
    t := RlGetProductInShop(product_id, shop_code)
    if t.Quantity > 0 {
      res.Quantity += t.Quantity
      res.Shops = append(res.Shops, (*GetShopByID(shop_code)))
    }
  }
  for _, warehouse_code := range GetWarehousesByRegionID(region_id) {
    t := RlGetProductInWarehouse(product_id, warehouse_code)
    if t.Quantity > 0 {
      res.Quantity += t.Quantity
      res.Warehouses = append(res.Warehouses, (*GetWarehouseByCode(warehouse_code)))
    }
  }
  return res
}
