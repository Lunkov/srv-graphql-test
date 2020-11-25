package main

import (
  "fmt"
  "math/rand"
  "strconv"
  "testing"
  "github.com/stretchr/testify/assert"
)

func BenchmarkWrite(b *testing.B) {
  fillData4Tests()
  //LoadAll()
  b.ResetTimer()
  SaveAll()
}

func BenchmarkRead(b *testing.B) {
  LoadAll()
}

func BenchmarkWHSerial(b *testing.B) {
  //fillData4Tests()
  LoadAll()
  
  //assert.Equal(b, &Warehouse{CODE:"a4", Name:"Name_WH_4", Description:""},   GetWarehouseByCode("a4"))
  //assert.Equal(b, &Warehouse{CODE:"a23", Name:"Name_WH_23", Description:""}, GetWarehouseByCode("a23"))

  //assert.Equal(b, &Shop{CODE:"s1",   Name:"SHOP_1",   Description:""}, GetShopByID("s1"))
  //assert.Equal(b, &Shop{CODE:"s543", Name:"SHOP_543", Description:""}, GetShopByID("s543"))

  assert.Equal(b, &Product{ID:1,      Name:"Product_1", Description:""}, GetProductByID(1))
  assert.Equal(b, &Product{ID:543,    Name:"Product_543", Description:""}, GetProductByID(543))
  id := ProductCount() - 2
  assert.Equal(b, &Product{ID:id, Name:fmt.Sprintf("Product_%d", id), Description:""}, GetProductByID(id))
  
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    Product_id := rand.Int63n(ProductCount())
    p := GetProduct4Sale(Product_id, rand.Int63n(RegionCount()))
    
    assert.Equal(b, GetProductByID(Product_id), &p.Prod)
  }

}


func BenchmarkWHParallel(b *testing.B) {
  //fillData4Tests()
  LoadAll()
  
  assert.Equal(b, &Product{ID:1,      Name:"Product_1", Description:""}, GetProductByID(1))
  assert.Equal(b, &Product{ID:543,    Name:"Product_543", Description:""}, GetProductByID(543))
  id := ProductCount() - 2
  assert.Equal(b, &Product{ID:id, Name:fmt.Sprintf("Product_%d", id), Description:""}, GetProductByID(id))
  
  b.ResetTimer()
  for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
      b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
          product_id := rand.Int63n(ProductCount())
          p := GetProduct4Sale(product_id, rand.Int63n(RegionCount()))
          
          assert.Equal(b, GetProductByID(product_id), &p.Prod)
        }
      })
    })
  }
}
