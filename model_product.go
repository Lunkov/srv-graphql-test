package main

import (
  "sync"
  "encoding/gob"
  "os"
  "unsafe"
  "github.com/golang/glog"
  "strconv"
  "github.com/graphql-go/graphql"
)

// Товары

type Product struct {
  ID           int64     `db:"id"             json:"id"              yaml:"id"`
  Name         string    `db:"name"           json:"name"            yaml:"name"`
  Description  string    `db:"description"    json:"description"     yaml:"description"`
}


var productType = graphql.NewObject(
   graphql.ObjectConfig{
      Name: "Product",
      Fields: graphql.Fields{
         "id": &graphql.Field{
            Type: graphql.String,
         },
         "name": &graphql.Field{
            Type: graphql.String,
         },
         "description": &graphql.Field{
            Type: graphql.String,
         },
      },
   },
)

var memG map[int64]Product
var muG   sync.RWMutex

func ProductInit(max int64) {
  memG = make(map[int64]Product, max)
  glog.Infof("LOG: Products: max          = %d", max)
  glog.Infof("LOG: Products: sizeof(item) = %d", unsafe.Sizeof(Product{}))
  glog.Infof("LOG: Products: sizeof(map)  = %d", unsafe.Sizeof(memG))
}

func ProductCount() int64 {
  return int64(len(memG))
}

func GetProducts(offset int, limit int) []Product {
  res := make([]Product, limit)
  i := 0
  for _, item := range memG {
    if offset <= i &&  offset + limit > i {
      res = append(res, item)
    }
    if offset + limit < i {
      break
    }
    i++
  }
  glog.Infof("LOG: Products %v", res)
  return res
}

func ProductAppend(info *Product) {
  muG.Lock()
  memG[info.ID] = *info
  muG.Unlock()
}

func GetProductByID(id int64) (*Product) {
  muG.RLock()
  item, ok := memG[id]
  muG.RUnlock()
  if ok {
    glog.Infof("LOG: Product(%d) = %v", id, item)
    return &item
  }
  return nil
}

func WriteFileProducts(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, _ := os.Create(filename)
  defer file.Close()
  encoder := gob.NewEncoder(file)
  encoder.Encode(memG)
}

func LoadFileProducts(wg *sync.WaitGroup, filename string) {
  defer wg.Done()
  file, err := os.Open(filename)
  if err !=nil {
    glog.Errorf("ERR: Load(%s): %v", filename, err)
    return
  }
  defer file.Close()
  
  decoder := gob.NewDecoder(file)
  err = decoder.Decode(&memG)
  if err != nil {
    glog.Errorf("ERR: Decoder(%s): %v", filename, err)
    return
  }
}


func ProductGQL() {
  AppendFields2GraphQL("product", &graphql.Field{
			Type: productType,
      Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                  Description: "id of the products",
                  Type:graphql.NewNonNull(graphql.String),
                },
              },
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        id, err := strconv.Atoi(p.Args["id"].(string))
        if err != nil {
          return nil, err
        }
				return GetProductByID(int64(id)), nil
			},
		})
    
	AppendFields2GraphQL("products", &graphql.Field{
			Type: graphql.NewList(productType),
      Args: graphql.FieldConfigArgument{
               "offset": &graphql.ArgumentConfig{
                  Type: graphql.Int,
               },
               "limit": &graphql.ArgumentConfig{
                  Type: graphql.Int,
               },
            },      
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
        offset, ok := p.Args["offset"].(int)
        if !ok {
          offset = 0
        }
        limit, ok := p.Args["limit"].(int)
        if !ok {
          limit = 1000
        }
        glog.Infof("LOG: Products(offset=%d, limit=%d)", offset, limit)
				return GetProducts(offset, limit), nil
			},
    })
}
